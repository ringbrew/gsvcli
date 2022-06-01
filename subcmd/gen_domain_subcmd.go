package subcmd

import (
	"errors"
	"log"
)

type GenDomainSubCmd struct {
	Domain string
}

func (sc *GenDomainSubCmd) Parse(args []string) error {
	if len(args) == 0 {
		return errors.New("invalid project name")
	}
	sc.Domain = args[0]
	return nil
}

func (sc *GenDomainSubCmd) Process() error {
	log.Printf("[INFO] gen domain command is running\n")

	p := NewGenDomain(sc.Domain)
	return p.Process()
}
