package utils

import (
	"os"
	"os/user"
)

func GetPath(configPath string, exchange string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	if dir == "/root" || dir == "/home" {
		return "/config/" + configPath + exchange + ".json"
	}
	return os.Getenv("GOPATH") + "/src/github.com/diadata-org/lumina-library/config/" + configPath + exchange + ".json"
}
