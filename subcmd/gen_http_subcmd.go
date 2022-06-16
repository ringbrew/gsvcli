package subcmd

import (
	"errors"
	"log"
)

type GenHttpSubCmd struct {
	Domain string
}

func (sc *GenHttpSubCmd) Parse(args []string) error {
	if len(args) == 0 {
		return errors.New("invalid domain")
	}
	sc.Domain = args[0]
	return nil
}

func (sc *GenHttpSubCmd) Process() error {
	log.Printf("[INFO] gen http command is running\n")

	p := NewGenHttp(sc.Domain)
	return p.Process()
}
