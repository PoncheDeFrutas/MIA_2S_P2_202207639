package commands

import (
	"backend/global"
	"backend/structures"
	"backend/utils"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type MkFile struct {
	Path string
	R    bool
	Size int
	Cont string
}

func ParserMkFile(tokens []string) (string, error) {
	cmd := &MkFile{}

	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`(?i)-path(?-i)="[^"]+"|(?i)-path(?-i)=\S+|(?i)-r|(?i)-size(?-i)=\d+|(?i)-cont(?-i)="[^"]+"|(?i)-cont(?-i)=\S+`)
	matches := re.FindAllString(args, -1)

	for _, match := range matches {
		var key, value string
		var err error

		if match == "-r" {
			key = "-r"
			value = ""
		} else {
			key, value, err = utils.ParseToken(match)
			if err != nil {
				return "", err
			}

			if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
				value = strings.Trim(value, "\"")
			}
		}

		switch key {
		case "-path":
			if value == "" {
				return "", fmt.Errorf("invalid path: %s", value)
			}
			cmd.Path = value
		case "-r":
			cmd.R = true
		case "-size":
			num, err := strconv.Atoi(value)
			if err != nil || num < 0 {
				return "", fmt.Errorf("invalid size: %d", num)
			}
			cmd.Size = num
		case "-cont":
			if value == "" {
				return "", fmt.Errorf("invalid content: %s", value)
			}
			cmd.Cont = value
		}
	}

	if cmd.Path == "" {
		return "", fmt.Errorf("path is required")
	}

	if err := cmd.commandMkFile(); err != nil {
		return "", err
	}

	return cmd.Print(), nil
}

func (cmd *MkFile) commandMkFile() error {
	_, _, err := global.GetLoggedUser()
	if err != nil {
		return fmt.Errorf("you must be logged in")
	}

	mountedPartition, partitionPath, err := global.GetMountedPartition(global.LoggedPartition)
	if err != nil {
		return err
	}

	sb := &structures.SuperBlock{}
	if err := sb.ReadSuperBlock(partitionPath, int64(mountedPartition.PartStart)); err != nil {
		return err
	}

	array := strings.Split(cmd.Path, "/")
	var result []string
	for _, part := range array {
		if part != "" {
			result = append(result, part)
		}
	}

	if err := sb.CreateNewInode(partitionPath, result, 0, true, cmd.R); err != nil {
		return err
	}

	cont := generateNumberString(cmd.Size)
	if _, err := sb.WriteFile(partitionPath, int32(0), result, cont); err != nil {
		return err
	}

	if cmd.Cont != "" {
		content, err := ioutil.ReadFile(cmd.Cont)
		if err != nil {
			log.Fatalf("Error al leer el archivo: %v", err)
		}

		fileContent := string(content)
		if _, err := sb.WriteFile(partitionPath, int32(0), result, fileContent); err != nil {
			return err
		}
	}

	if err := sb.WriteSuperBlock(partitionPath, int64(mountedPartition.PartStart), int64(mountedPartition.PartStart+int32(binary.Size(sb)))); err != nil {
		return err
	}

	return nil
}

func generateNumberString(n int) string {
	base := "0123456789"

	repeatCount := n / len(base)
	remainder := n % len(base)

	return strings.Repeat(base, repeatCount) + base[:remainder]
}

func (cmd *MkFile) Print() string {
	return fmt.Sprintf("File created successfully in %s", cmd.Path)
}
