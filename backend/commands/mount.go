package commands

import (
	"backend/global"
	"backend/structures"
	"backend/utils"
	"fmt"
	"regexp"
	"strings"
)

type Mount struct {
	Path string
	Name string
}

func ParserMount(tokens []string) (string, error) {
	cmd := &Mount{}

	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`(?i)-path(?-i)="[^"]+"|(?i)-path(?-i)=\S+|(?i)-name(?-i)="[^"]+"|(?i)-name(?-i)=\S+`)
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
		case "-path":
			if value == "" {
				return "", fmt.Errorf("invalid path: %s", value)
			}
			cmd.Path = value
		case "-name":
			if value == "" {
				return "", fmt.Errorf("invalid name: %s", value)
			}
			cmd.Name = value
		}
	}

	if cmd.Path == "" {
		return "", fmt.Errorf("missing path")
	}

	if cmd.Name == "" {
		return "", fmt.Errorf("missing name")
	}

	if result, err := cmd.commandMount(); err != nil {
		return "", err
	} else if result != "" {
		return result, nil
	}

	return "MOUNT FAIL", nil
}

func (cmd *Mount) commandMount() (string, error) {
	mbr := structures.MBR{}

	if err := mbr.ReadMBR(cmd.Path); err != nil {
		return "", err
	}

	partition, indexPartition := mbr.GetPartitionByName(cmd.Name)

	if partition == nil {
		return "", fmt.Errorf("partition not found")
	}

	if partition.PartType != 'P' {
		return "", fmt.Errorf("partition is not primary")
	}

	idPartition, err := cmd.GenerateIdPartition(indexPartition)

	if err != nil {
		return "", err
	}

	global.MountedPartitions[idPartition] = cmd.Path

	if err := partition.MountPartition(indexPartition, idPartition); err != nil {
		return "", err
	}

	mbr.MbrPartition[indexPartition] = *partition

	if err := mbr.WriteMBR(cmd.Path); err != nil {
		return "", err
	}

	return fmt.Sprintf("partition mounted successfully with id: %s", idPartition), nil
}

func (cmd *Mount) GenerateIdPartition(indexPartition int) (string, error) {
	letter, err := utils.GetLetter(cmd.Path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%d%s", global.Carnet, indexPartition+1, letter), nil
}
