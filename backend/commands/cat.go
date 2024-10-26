package commands

import (
	"backend/global"
	"backend/structures"
	"backend/utils"
	"fmt"
	"regexp"
	"strings"
)

type Cat struct {
	FileN []string
}

func ParserCat(tokens []string) (string, error) {
	cmd := &Cat{}

	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`(?i)-file\d+(?-i)="[^"]+"|(?i)-file\d+(?-i)=\S+`)
	matches := re.FindAllString(args, -1)

	for _, match := range matches {
		key, value, err := utils.ParseToken(match)
		if err != nil {
			return "", err
		}

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		if strings.HasPrefix(key, "-file") {
			cmd.FileN = append(cmd.FileN, value)
		}
	}

	if len(cmd.FileN) == 0 {
		return "", nil
	}

	return cmd.commandCat()

}

func (cmd *Cat) commandCat() (string, error) {
	if !global.IsUserLogged() {
		return "", fmt.Errorf("you must be logged")
	}

	mountedPartition, partitionPath, err := global.GetMountedPartition(global.LoggedPartition)
	if err != nil {
		return "", err
	}

	sb := &structures.SuperBlock{}
	if err := sb.ReadSuperBlock(partitionPath, int64(mountedPartition.PartStart)); err != nil {
		return "", err
	}

	var sb2 strings.Builder

	for _, file := range cmd.FileN {
		array := strings.Split(file, "/")
		var result []string
		for _, part := range array {
			if part != "" {
				result = append(result, part)
			}
		}
		response := sb.GetFile(partitionPath, 0, result)
		if response != "" {
			sb2.WriteString(response)
			sb2.WriteString("\n")
		}
	}

	return sb2.String(), nil
}
