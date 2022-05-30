package subcmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GenGrpc struct {
	cache     string
	protoRepo string
	module    string
}

func NewGenGrpc() *GenGrpc {
	tmpDir := os.TempDir()
	cache := filepath.Join(tmpDir, "proto_dep")

	module := ""
	if file, err := os.Open("go.mod"); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()
			if strings.HasPrefix(text, "module ") {
				module = strings.TrimPrefix(text, "module ")
				break
			}
		}
	} else {
		log.Fatal(err.Error())
	}

	if module == "" {
		log.Fatal("please init your go module with go mod init")
	}

	return &GenGrpc{
		cache:     cache,
		protoRepo: "proto-center",
		module:    module,
	}
}

func (g GenGrpc) Process(protoFile string) error {
	if err := g.Setup(); err != nil {
		return err
	}

	if err := g.Gen(protoFile, "internal/delivery"); err != nil {
		return err
	}

	if err := g.ModTidy(); err != nil {
		return err
	}

	if err := GoFmt(); err != nil {
		return err
	}

	return nil
}

/*
	protoc -I ../proto/ ../proto/sample/user_demo.proto --go_out=plugins=grpc,paths=import:./external
*/
func (g GenGrpc) Setup() error {
	if err := os.RemoveAll(g.cache); err != nil {
		return err
	}

	if err := g.CloneProtoRepo(); err != nil {
		return err
	}

	if err := g.SetGoEnv(); err != nil {
		return err
	}

	if err := g.PullProtoRepo(); err != nil {
		return err
	}

	return nil
}

func (g GenGrpc) SetGoEnv() error {
	c := exec.Command("go", "env", `GOPROXY="https://goproxy.cn"`)
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	c = exec.Command("go", "env", `GO111MODULE="on"`)
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	return nil
}

func (g GenGrpc) CloneProtoRepo() error {
	c := exec.Command("git", "clone", fmt.Sprintf("https://github.com/ringbrew/%s.git", g.protoRepo), filepath.Base(g.cache))
	c.Dir = filepath.Dir(g.cache)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}
	return nil
}

func (g GenGrpc) PullProtoRepo() error {
	c := exec.Command("git", "checkout", "master")
	c.Dir = g.cache
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	c = exec.Command("git", "fetch", "origin", "master")
	c.Dir = g.cache
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	c = exec.Command("git", "reset", "origin/master", "--hard")
	c.Dir = g.cache
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	return nil
}

/*
	protoc -I ../proto/ ../proto/sample/user_demo.proto --go_out=plugins=grpc,paths=import:./external
*/
func (g GenGrpc) Gen(serviceProto, delivery string) error {
	if err := os.MkdirAll("export", os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll("openapi", os.ModePerm); err != nil {
		return err
	}

	if err := g.protoC(); err != nil {
		return err
	}

	return nil
}

func (g GenGrpc) protoC() error {
	domainPath := filepath.Join(g.cache, "proto")

	if !Exists(domainPath) {
		if err := os.MkdirAll(domainPath, os.ModePerm); err != nil {
			return err
		}
	}

	pf, err := CopyDir("proto", domainPath)
	if err != nil {
		return err
	}

	protoFile := make([]string, 0, len(pf))
	for _, v := range pf {
		protoFile = append(protoFile, filepath.Join(domainPath, v))
	}

	args := make([]string, 0, 20)
	args = append(args, "protoc")
	args = append(args, "-I", g.cache,
		fmt.Sprintf("--go_out=./"),
		fmt.Sprintf("--go_opt=module=%s", g.module),
		fmt.Sprintf("--go-grpc_out=./"),
		fmt.Sprintf("--go-grpc_opt=module=%s", g.module),
		fmt.Sprintf("--grpc-gateway_out=:./"),
		fmt.Sprintf("--grpc-gateway_opt=logtostderr=true"),
		fmt.Sprintf("--grpc-gateway_opt=module=%s", g.module),
		fmt.Sprintf("--openapiv2_out=./openapi"),
		fmt.Sprintf("--openapiv2_opt=logtostderr=true"),
		fmt.Sprintf("--openapiv2_opt=allow_merge=true"),
		fmt.Sprintf("--openapiv2_opt=merge_file_name=%s", filepath.Base(g.module)),
	)

	args = append(args, protoFile...)

	log.Println(args)
	c := exec.Command(args[0], args[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}

func (g GenGrpc) ModTidy() error {
	c := exec.Command("go", "mod", "tidy")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}
	return nil
}
