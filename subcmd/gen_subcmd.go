package subcmd

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type GenSubCmd struct {
	Args map[string]string
}

func (sc *GenSubCmd) Parse(args []string) error {
	if len(args) == 0 {
		return errors.New("invalid gen args")
	}

	foundProtoFile := false

	sc.Args = make(map[string]string)

	for i := 0; i < len(args); i++ {
		if !strings.HasPrefix(args[i], "-") {
			if foundProtoFile {
				return errors.New("detect two proto file param")
			} else {
				foundProtoFile = true
				sc.Args["file"] = args[i]
			}
		} else {
			a := strings.Split(args[i], "=")
			if len(a) != 2 {
				return fmt.Errorf("invalid param[%s]", args[i])
			}

			sc.Args[a[0][1:]] = a[1]
		}
	}

	return nil
}

func (sc *GenSubCmd) Process() error {
	log.Printf("gen command is running[%v]\n", sc.Args)

	if sc.Args["file"] == "" {
		return fmt.Errorf("invalid file input")
	}

	p := NewGenGrpc()
	return p.Process(sc.Args["file"])
}
