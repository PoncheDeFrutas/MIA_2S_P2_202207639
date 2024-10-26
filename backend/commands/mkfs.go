package commands

import (
	"backend/global"
	"backend/structures"
	"backend/utils"
	"encoding/binary"
	"fmt"
	"regexp"
	"strings"
)

type MkFs struct {
	Id   string
	Type string
}

func ParserMkFs(tokens []string) (string, error) {
	cmd := &MkFs{}

	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`(?i)-id(?-i)=\S+|(?i)-type(?-i)=\S+`)
	matches := re.FindAllString(args, -1)

	for _, match := range matches {
		key, value, err := utils.ParseToken(match)
		if err != nil {
			return "", err
		}

		switch key {
		case "-id":
			if value == "" {
				return "", fmt.Errorf("invalid id: %s", value)
			}
			cmd.Id = value
		case "-type":
			if strings.ToLower(value) != "full" {
				return "", fmt.Errorf("invalid type: %s", value)
			}
			cmd.Type = value
		}
	}

	if cmd.Id == "" {
		return "", fmt.Errorf("missing id")
	}

	if cmd.Type == "" {
		cmd.Type = "full"
	}

	if err := cmd.commandMkFs(); err != nil {
		return "", err
	}

	return cmd.Print(), nil
}

func (cmd *MkFs) commandMkFs() error {
	mountedPartition, partitionPath, err := global.GetMountedPartition(cmd.Id)
	if err != nil {
		return err
	}

	// TODO verify if the partition is mounted
	fmt.Println("MOUNTED PARTITION: ")
	mountedPartition.Print()

	n := mountedPartition.CalculateN()
	fmt.Println("N: ", n)

	superBlock := structures.SuperBlock{}
	superBlock.CreateSuperBlock(mountedPartition.PartStart, n)

	fmt.Println("SUPER BLOCK: ")
	superBlock.Print()

	if err := superBlock.CreateBitMaps(partitionPath); err != nil {
		return err
	}

	if err := superBlock.CreateUserFile(partitionPath); err != nil {
		return err
	}

	superBlock.Print()

	if err := superBlock.WriteSuperBlock(partitionPath, int64(mountedPartition.PartStart), int64(mountedPartition.PartStart+int32(binary.Size(superBlock)))); err != nil {
		return err
	}

	return nil
}

func (cmd *MkFs) Print() string {
	return fmt.Sprintf("File system created successfully in partition %s", cmd.Id)
}
