package subcmd

import (
	"log"
)

type VersionSubCmd struct {
}

func (sc *VersionSubCmd) Parse(args []string) error {
	return nil
}

func (sc *VersionSubCmd) Process() error {
	log.Println("[INFO] version 1.0.0")
	return nil
}
