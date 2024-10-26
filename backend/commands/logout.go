package commands

import (
	"backend/global"
)

func ParserLogout(tokens []string) (string, error) {
	if text, err := global.LogUserOut(); err != nil {
		return "", err
	} else {
		return "logout" + text, nil
	}
}
