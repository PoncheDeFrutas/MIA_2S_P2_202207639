package commands

import (
	"backend/structures"
	"backend/utils"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type MkDisk struct {
	Size int
	Fit  string
	Unit string
	Path string
}

func ParserMkDisk(tokens []string) (string, error) {
	cmd := &MkDisk{}

	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`(?i)-size(?-i)=\d+|(?i)-unit(?-i)=[kKmM]|(?i)-fit(?-i)=[bBfFwW]{2}|(?i)-path(?-i)="[^"]+"|(?i)-path(?-i)=\S+|-.+`)
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
		case "-fit":
			value = strings.ToUpper(value)
			if value != "BF" && value != "FF" && value != "WF" {
				return "", fmt.Errorf("invalid fit: %s", value)
			}
			cmd.Fit = value
		case "-unit":
			value = strings.ToUpper(value)
			if value != "K" && value != "M" {
				return "", fmt.Errorf("invalid unit: %s", value)
			}
			cmd.Unit = value
		case "-path":
			if value == "" {
				return "", fmt.Errorf("invalid path: %s", value)
			}
			cmd.Path = value
		default:
			return "", fmt.Errorf("unknown parameter: %s", key)
		}
	}

	if cmd.Size == 0 {
		return "", fmt.Errorf("missing size")
	}
	if cmd.Fit == "" {
		cmd.Fit = "FF"
	}
	if cmd.Unit == "" {
		cmd.Unit = "M"
	}
	if cmd.Path == "" {
		return "", fmt.Errorf("missing path")
	}

	if err := cmd.commandMkDisk(); err != nil {
		return "", err
	}

	return cmd.Print(), nil
}

func (cmd *MkDisk) commandMkDisk() error {
	sizeInBytes, err := utils.ConvertToBytes(cmd.Size, cmd.Unit)
	if err != nil {
		return err
	}

	if err := cmd.createDisk(sizeInBytes); err != nil {
		return err
	}

	mbr := &structures.MBR{}

	if err := mbr.CreateMBR(sizeInBytes, cmd.Fit); err != nil {
		return err
	}

	if err := mbr.WriteMBR(cmd.Path); err != nil {
		return err
	}

	mbr.Print()

	return nil
}

func (cmd *MkDisk) createDisk(sizeInBytes int) error {
	if err := os.MkdirAll(filepath.Dir(cmd.Path), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(cmd.Path)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	buffer := make([]byte, 1024*1024)

	for sizeInBytes > 0 {
		writeSize := len(buffer)

		if sizeInBytes < len(buffer) {
			writeSize = sizeInBytes
		}

		if _, err := file.Write(buffer[:writeSize]); err != nil {
			return err
		}

		sizeInBytes -= writeSize
	}

	return nil
}

func (cmd *MkDisk) Print() string {
	return fmt.Sprintf("Disk created successfully at: %s\nSize: %d%s\nFit: %s", cmd.Path, cmd.Size, cmd.Unit, cmd.Fit)
}
