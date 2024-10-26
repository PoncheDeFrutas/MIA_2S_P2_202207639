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

type RmUSR struct {
	User string
}

func ParserRmUSR(tokens []string) (string, error) {
	cmd := &RmUSR{}

	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`(?i)-user(?-i)="[^"]+"|(?i)-user(?-i)=\S+`)
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
		case "-user":
			if value == "" {
				return "", fmt.Errorf("invalid name: %s", value)
			}
			cmd.User = value
		}
	}

	if cmd.User == "" {
		return "", fmt.Errorf("name is required")
	}

	if err := cmd.commandRmGRP(); err != nil {
		return "", err
	}

	return cmd.Print(), nil
}

func (cmd *RmUSR) commandRmGRP() error {
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

	if err := global.RemoveUser(cmd.User); err != nil {
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

func (cmd *RmUSR) Print() string {
	return fmt.Sprintf("User %s removed", cmd.User)
}
