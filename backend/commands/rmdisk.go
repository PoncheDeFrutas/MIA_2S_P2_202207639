package commands

import (
	common "backend/utils"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type RmDisk struct {
	Path string
}

func ParserRmDisk(tokens []string) (string, error) {
	cmd := &RmDisk{}

	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`(?i)-path(?-i)="[^"]+"|(?i)-path(?-i)=\S+`)
	matches := re.FindAllString(args, -1)

	for _, match := range matches {
		key, value, err := common.ParseToken(match)
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
		}
	}

	if cmd.Path == "" {
		return "", fmt.Errorf("missing path")
	}

	err := commandRmDisk(cmd.Path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("RmDisk: %s", cmd.Path), nil
}

func commandRmDisk(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}
