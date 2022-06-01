package subcmd

import (
	"log"
)

type GenGrpcSubCmd struct {
	Path string
	//Args map[string]string
}

func (sc *GenGrpcSubCmd) Parse(args []string) error {
	if len(args) > 0 {
		sc.Path = args[0]
	} else {
		sc.Path = "proto"
	}

	//if len(args) == 0 {
	//	return errors.New("invalid gen args")
	//}

	//
	//foundProtoPath := false
	//
	//sc.Args = make(map[string]string)
	//
	//for i := 0; i < len(args); i++ {
	//	if !strings.HasPrefix(args[i], "-") {
	//		if foundProtoPath {
	//			return errors.New("detect two proto file param")
	//		} else {
	//			foundProtoPath = true
	//			sc.Args["path"] = args[i]
	//		}
	//	} else {
	//		a := strings.Split(args[i], "=")
	//		if len(a) != 2 {
	//			return fmt.Errorf("invalid param[%s]", args[i])
	//		}
	//
	//		sc.Args[a[0][1:]] = a[1]
	//	}
	//}

	return nil
}

func (sc *GenGrpcSubCmd) Process() error {
	log.Printf("gen grpc command is running[%v]\n", sc.Path)
	return NewGenGrpc().Process(sc.Path)
}
