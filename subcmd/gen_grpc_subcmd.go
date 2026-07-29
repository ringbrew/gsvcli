package subcmd

import (
	"errors"
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
	importPath, protoPath, err := parseGrpcArgs(args)
	if err != nil {
		return err
	}

	sc.ImportPath = importPath
	sc.ProtoPath = protoPath
	return nil
}

// parseGrpcArgs 解析 -I=importPath / -P=protoPath。命令关掉了 cobra 的 flag 解析
// (proto 路径里可能带 cobra 会误读的字符)，所以这里自己解，两处入口共用一份。
//
// importPath 为空时回落到 GOPROTO：第三方依赖 proto(google/api、openapiv2 等)
// 集中放一个目录，用环境变量指过去，各仓就不用各自维护依赖路径。
// protoPath 为空时默认 proto，即每个项目固定的 proto 目录。
func parseGrpcArgs(args []string) (importPath string, protoPath string, err error) {
	argsMap := make(map[string]string)

	for i := 0; i < len(args); i++ {
		if args[i] == "-h" || args[i] == "--help" {
			return "", "", errHelp
		}

		a := strings.SplitN(args[i], "=", 2)
		if len(a) != 2 || len(a[0]) < 2 || a[0][0] != '-' {
			return "", "", fmt.Errorf("invalid param[%s], usage: gsv gen grpc [-I=importPath] [-P=protoPath]", args[i])
		}

		argsMap[a[0][1:]] = a[1]
	}

	importPath = argsMap["I"]
	if importPath == "" {
		importPath = os.Getenv("GOPROTO")
	}

	protoPath = argsMap["P"]
	if protoPath == "" {
		protoPath = "proto"
	}

	return importPath, protoPath, nil
}

// errHelp 表示用户要的是用法说明，不是一次失败的生成。
var errHelp = errors.New("help requested")

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
			importPath, protoPath, err := parseGrpcArgs(args)
			if err != nil {
				if errors.Is(err, errHelp) {
					_ = cmd.Usage()
					return
				}
				log.Fatal(err.Error())
			}

			log.Printf("gen grpc command is running -I[%s], -P[%s]\n", importPath, protoPath)

			if err := NewGenGrpc().Process(importPath, protoPath); err != nil {
				log.Fatal(err.Error())
			}
		},
	}
}
