package subcmd

import (
	"fmt"
	"log"
	"strings"
)

type GenGrpcSubCmd struct {
	ImportPath string
	ProtoPath  string
	//Args map[string]string
}

func (sc *GenGrpcSubCmd) Parse(args []string) error {
	//if len(args) > 0 {
	//	sc.Path = args[0]
	//} else {
	//	sc.Path = "proto"
	//}

	//if len(args) == 0 {
	//	return errors.New("invalid gen args")
	//}

	argsMap := make(map[string]string)

	for i := 0; i < len(args); i++ {

		a := strings.Split(args[i], "=")
		if len(a) != 2 {
			return fmt.Errorf("invalid param[%s]", args[i])
		}

		argsMap[a[0][1:]] = a[1]
	}

	sc.ImportPath = argsMap["I"]
	sc.ProtoPath = argsMap["P"]

	if sc.ProtoPath == "" {
		sc.ProtoPath = "proto"
	}
	return nil
}

func (sc *GenGrpcSubCmd) Process() error {
	log.Printf("gen grpc command is running -I[%s], -P[%s]\n", sc.ImportPath, sc.ProtoPath)
	return NewGenGrpc().Process(sc.ImportPath, sc.ProtoPath)
}
