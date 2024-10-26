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

type MkUSR struct {
	User string
	Pass string
	Grp  string
}

func ParserMkUSR(tokens []string) (string, error) {
	cmd := &MkUSR{}

	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`(?i)-user(?-i)="[^"]+"|(?i)-user(?-i)=\S+|(?i)-pass(?-i)="[^"]+"|(?i)-pass(?-i)=\S+|(?i)-grp(?-i)="[^"]+"|(?i)-grp(?-i)=\S+`)
	matches := re.FindAllString(args, -1)

	for _, match := range matches {
		key, value, err := utils.ParseToken(match)
		if err != nil {
			return "", err
		}

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		if len(value) > 10 {
			return "", fmt.Errorf("%s must be at most 10 characters long", key)
		}

		switch key {
		case "-user":
			if value == "" {
				return "", fmt.Errorf("invalid user: %s", value)
			}
			cmd.User = value
		case "-pass":
			if value == "" {
				return "", fmt.Errorf("invalid pass: %s", value)
			}
			cmd.Pass = value
		case "-grp":
			if value == "" {
				return "", fmt.Errorf("invalid group: %s", value)
			}
			cmd.Grp = value
		}
	}

	if cmd.User == "" {
		return "", fmt.Errorf("user is required")
	}

	if cmd.Pass == "" {
		return "", fmt.Errorf("password is required")
	}

	if cmd.Grp == "" {
		return "", fmt.Errorf("group is required")
	}

	if err := cmd.commandMkUSR(); err != nil {
		return "", err
	}

	return cmd.Print(), nil
}

func (cmd *MkUSR) commandMkUSR() error {
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

	if err := global.AddUserToGroup(cmd.User, cmd.Pass, cmd.Grp); err != nil {
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

func (cmd *MkUSR) Print() string {
	return fmt.Sprintf("user %s created in group %s", cmd.User, cmd.Grp)
}
