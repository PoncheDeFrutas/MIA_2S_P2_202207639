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

type RmGRP struct {
	Name string
}

func ParserRmGRP(tokens []string) (string, error) {
	cmd := &RmGRP{}

	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`(?i)-name(?-i)="[^"]+"|(?i)-name(?-i)=\S+`)
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
		case "-name":
			if value == "" {
				return "", fmt.Errorf("invalid name: %s", value)
			}
			cmd.Name = value
		}
	}

	if cmd.Name == "" {
		return "", fmt.Errorf("name is required")
	}

	if err := cmd.commandRmGRP(); err != nil {
		return "", err
	}

	return cmd.Print(), nil
}

func (cmd *RmGRP) commandRmGRP() error {
	if user, _, err := global.GetLoggedUser(); user != "root" || err != nil {
		return fmt.Errorf("permission denied")
	}

	mountedPartition, partitionPath, err := global.GetMountedPartition(global.LoggedPartition)
	if err != nil {
		return err
	}

	sb := &structures.SuperBlock{}
	if err := sb.ReadSuperBlock(partitionPath, int64(mountedPartition.PartStart)); err != nil {
		return err
	}

	array := []string{"users.txt"}

	texto := sb.GetFile(partitionPath, 0, array)

	global.ParserUserData(texto)

	if err := global.RemoveGroup(cmd.Name); err != nil {
		return err
	}

	if _, err := sb.WriteFile(partitionPath, int32(0), array, global.ConvertToString()); err != nil {
		return err
	}

	if err := sb.WriteSuperBlock(partitionPath, int64(mountedPartition.PartStart), int64(mountedPartition.PartStart+int32(binary.Size(sb)))); err != nil {
		return err
	}

	return nil
}

func (cmd *RmGRP) Print() string {
	return fmt.Sprintf("group %s removed", cmd.Name)
}
