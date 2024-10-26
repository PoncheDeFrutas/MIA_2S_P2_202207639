package commands

import (
	"backend/structures"
	"backend/utils"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type FDisk struct {
	Size int
	Unit string
	Path string
	Type string
	Fit  string
	Name string
}

func ParserFDisk(tokens []string) (string, error) {
	cmd := &FDisk{}

	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`(?i)-size(?-i)=\d+|(?i)-unit(?-i)=[bBkKmM]|(?i)-fit(?-i)=[bBfF]{2}|(?i)-path(?-i)="[^"]+"|(?i)-path(?-i)=\S+|(?i)-type(?-i)=[pPeElL]|(?i)-name(?-i)="[^"]+"|(?i)-name(?-i)=\S+`)
	matches := re.FindAllString(args, -1)

	for _, match := range matches {
		key, value, err := utils.ParseToken(match)
		if err != nil {
			return "", err
		}

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		switch key {
		case "-size":
			size, err := strconv.Atoi(value)
			if err != nil || size < 1 {
				return "", fmt.Errorf("invalid size: %s", value)
			}
			cmd.Size = size
		case "-unit":
			value = strings.ToUpper(value)
			if value != "B" && value != "K" && value != "M" {
				return "", fmt.Errorf("invalid unit: %s", value)
			}
			cmd.Unit = value
		case "-path":
			if value == "" {
				return "", fmt.Errorf("invalid path: %s", value)
			}
			cmd.Path = value
		case "-type":
			value = strings.ToUpper(value)
			if value != "P" && value != "E" && value != "L" {
				return "", fmt.Errorf("invalid type: %s", value)
			}
			cmd.Type = value
		case "-fit":
			value = strings.ToUpper(value)
			if value != "BF" && value != "FF" && value != "WF" {
				return "", fmt.Errorf("invalid fit: %s", value)
			}
			cmd.Fit = value
		case "-name":
			if value == "" {
				return "", fmt.Errorf("invalid name: %s", value)
			}
			cmd.Name = value
		default:
			return "", fmt.Errorf("unknown option: %s", key)
		}
	}

	if cmd.Size == 0 {
		return "", fmt.Errorf("missing size")
	}

	if cmd.Unit == "" {
		cmd.Unit = "K"
	}

	if cmd.Path == "" {
		return "", fmt.Errorf("missing path")
	}

	if cmd.Type == "" {
		cmd.Type = "P"
	}

	if cmd.Fit == "" {
		cmd.Fit = "WF"
	}

	if cmd.Name == "" {
		return "", fmt.Errorf("missing name")
	}

	if err := cmd.commandFDisk(); err != nil {
		return "", fmt.Errorf("commandFDisk failed: %w (cmd details: Size=%d, Unit=%s, Path=%s, Type=%s, Fit=%s, Name=%s)",
			err, cmd.Size, cmd.Unit, cmd.Path, cmd.Type, cmd.Fit, cmd.Name)
	}

	mbr := structures.MBR{}
	if err := mbr.ReadMBR(cmd.Path); err != nil {
		return "", err
	}
	return cmd.Print(), nil
}

func (cmd *FDisk) commandFDisk() error {
	sizeInBytes, err := utils.ConvertToBytes(cmd.Size, cmd.Unit)
	if err != nil {
		return err
	}

	mbr := &structures.MBR{}
	if err := mbr.ReadMBR(cmd.Path); err != nil {
		return err
	}

	switch cmd.Type {
	case "P":
		if err := cmd.createPrimaryPartition(mbr, sizeInBytes); err != nil {
			return err
		}
	case "E":
		if err := cmd.createExtendedPartition(mbr, sizeInBytes); err != nil {
			return err
		}
	case "L":
		if err := cmd.createLogicalPartition(mbr, sizeInBytes); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid type: %s", cmd.Type)
	}
	return nil
}

func (cmd *FDisk) createPrimaryPartition(mbr *structures.MBR, sizeInBytes int) error {
	indexPart, indexByte, err := cmd.findAvailableSpace(mbr, sizeInBytes)
	if err != nil {
		return err
	}

	if mbr.FreeNamePartition(cmd.Name) == false {
		return fmt.Errorf("name already exists: %s", cmd.Name)
	}

	partition := &mbr.MbrPartition[indexPart]
	partition.SetPartition(cmd.Type, cmd.Fit, indexByte, int32(sizeInBytes), cmd.Name)

	if err := mbr.WriteMBR(cmd.Path); err != nil {
		return err
	}

	return nil
}

func (cmd *FDisk) createExtendedPartition(mbr *structures.MBR, sizeInBytes int) error {
	if mbr.ExtendPartitionExist() {
		return fmt.Errorf("extended partition already exists")
	}

	indexPart, indexByte, err := cmd.findAvailableSpace(mbr, sizeInBytes)
	if err != nil {
		return err
	}

	if mbr.FreeNamePartition(cmd.Name) == false {
		return fmt.Errorf("name already exists: %s", cmd.Name)
	}

	partition := &mbr.MbrPartition[indexPart]
	partition.SetPartition(cmd.Type, cmd.Fit, indexByte, int32(sizeInBytes), cmd.Name)

	if err := mbr.WriteMBR(cmd.Path); err != nil {
		return err
	}

	ebr := &structures.EBR{}
	ebr.DefaultValue()

	if err := ebr.WriteEBR(cmd.Path, int64(partition.PartStart), int64(partition.PartSize+partition.PartStart)); err != nil {
		return err
	}

	return nil
}

func (cmd *FDisk) createLogicalPartition(mbr *structures.MBR, sizeInBytes int) error {
	if !mbr.ExtendPartitionExist() {
		return fmt.Errorf("extended partition does not exist")
	}

	start := int32(0)
	partition := mbr.GetExtendedPartition()
	start = partition.PartStart

	ebr := &structures.EBR{}
	if err := ebr.ReadEBR(cmd.Path, int64(partition.PartStart)); err != nil {
		return err
	}

	for ebr.PartNext != -1 {
		start = ebr.PartNext
		if err := ebr.ReadEBR(cmd.Path, int64(ebr.PartNext)); err != nil {
			return err
		}
	}

	ebr.SetEBR(cmd.Fit, start+30, int32(sizeInBytes), start+30+int32(sizeInBytes), cmd.Name)

	if err := ebr.WriteEBR(cmd.Path, int64(ebr.PartStart-30), int64(partition.PartStart+partition.PartSize)); err != nil {
		return err
	}

	def := &structures.EBR{}
	def.DefaultValue()

	if err := def.WriteEBR(cmd.Path, int64(ebr.PartNext), int64(partition.PartStart+partition.PartSize)); err != nil {
		return err
	}

	return nil
}

func (cmd *FDisk) findAvailableSpace(mbr *structures.MBR, sizeInBytes int) (int, int32, error) {
	indexPart := mbr.FindFreePartition()
	if indexPart == -1 {
		return -1, -1, fmt.Errorf("no free partition available")
	}

	objects := structures.ConvertToObjects(mbr.MbrPartition[:])
	indexByte := structures.FirstFit(objects, int32(sizeInBytes), int32(153), mbr.MbrSize)

	if indexByte == -1 {
		return -1, -1, fmt.Errorf("no space available for partition")
	}
	return indexPart, indexByte, nil
}

func (cmd *FDisk) Print() string {
	return fmt.Sprintf("FDISK\n Size: %d\n Unit: %s\n Path: %s\n Type: %s\n Fit: %s\n Name: %s\n", cmd.Size, cmd.Unit, cmd.Path, cmd.Type, cmd.Fit, cmd.Name)
}
