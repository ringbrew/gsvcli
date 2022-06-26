package subcmd

import (
	"os"
	"os/exec"
)

type InstallSubCmd struct {
}

func (sc *InstallSubCmd) Parse(args []string) error {
	return nil
}

func (sc *InstallSubCmd) Process() error {
	// install dependency
	dependency := []string{
		"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.10.1",
		"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.10.1",
		"github.com/ringbrew/protoc-gen-go-gsv@latest",
		"google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0",
		"github.com/ringbrew/protoc-go-inject-tag@v1.3.1",
	}

	/*
		go install \
		    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
			github.com/ringbrew/protoc-gen-go-gsv \
		    google.golang.org/protobuf/cmd/protoc-gen-go \
		    google.golang.org/grpc/cmd/protoc-gen-go-grpc \
	*/

	for _, v := range dependency {
		c := exec.Command("go", "install", "-v", v)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Dir = "/"
		if err := c.Run(); err != nil {
			return err
		}

	}
	return nil
}
