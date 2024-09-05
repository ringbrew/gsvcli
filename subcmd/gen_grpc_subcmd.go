package subcmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
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

func NewGrpcCommand() *cobra.Command {
	return &cobra.Command{
		Use:                "grpc",
		Short:              "generator for grpc code. usage: gsv gen grpc optional:-I=importPath optional:-P=protoPath",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			argsMap := make(map[string]string)

			for i := 0; i < len(args); i++ {
				a := strings.Split(args[i], "=")
				if len(a) != 2 {
					log.Fatal(fmt.Errorf("invalid param[%s]", args[i]).Error())
				}

				argsMap[a[0][1:]] = a[1]
			}

			importPath := argsMap["I"]
			protoPath := argsMap["P"]

			if protoPath == "" {
				protoPath = "proto"
			}

			if importPath == "" {
				importPath = os.Getenv("GOPROTO")
			}

			log.Printf("gen grpc command is running -I[%s], -P[%s]\n", importPath, protoPath)

			if err := NewGenGrpc().Process(importPath, protoPath); err != nil {
				log.Fatal(err.Error())
			}
		},
	}
}
