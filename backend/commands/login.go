package commands

import (
	"backend/global"
	"backend/structures"
	"backend/utils"
	"fmt"
	"regexp"
	"strings"
)

type Login struct {
	User string
	Pass string
	Id   string
}

func ParserLogin(tokens []string) (string, error) {
	cmd := &Login{}

	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`(?i)-user(?-i)="[^"]+"|(?i)-user(?-i)=\S+|(?i)-pass(?-i)="[^"]+"|(?i)-pass(?-i)=\S+|(?i)-id(?-i)="[^"]+"|(?i)-id(?-i)=\S+`)
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
				return "", fmt.Errorf("invalid user: %s", value)
			}
			cmd.User = value
		case "-pass":
			if value == "" {
				return "", fmt.Errorf("invalid pass: %s", value)
			}
			cmd.Pass = value
		case "-id":
			if value == "" {
				return "", fmt.Errorf("invalid id: %s", value)
			}
			cmd.Id = value
		default:
			return "", fmt.Errorf("unknown option: %s", key)
		}
	}

	if cmd.User == "" {
		return "", fmt.Errorf("user is required")
	}

	if cmd.Pass == "" {
		return "", fmt.Errorf("pass is required")
	}

	if cmd.Id == "" {
		return "", fmt.Errorf("id is required")
	}

	if err := cmd.commandLogin(); err != nil {
		return "", err
	}

	return cmd.Print(), nil
}

func (cmd *Login) commandLogin() error {
	if global.IsUserLogged() {
		return fmt.Errorf("a user is already logged")
	}

	mountedPartition, partitionPath, err := global.GetMountedPartition(cmd.Id)
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

	if err := global.LogUserIn(cmd.User, cmd.Pass, cmd.Id); err != nil {
		return err
	}

	return nil
}

func (cmd *Login) Print() string {
	return fmt.Sprintf("User %s logged in", cmd.User)
}
