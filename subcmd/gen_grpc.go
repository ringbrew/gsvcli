package subcmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
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
		log.Fatal("[FATAL] " + err.Error())
	}

	if module == "" {
		log.Fatal("[FATAL] please init your go module with go mod init")
	}

	return &GenGrpc{
		cache:     cache,
		protoRepo: "proto-center",
		module:    module,
	}
}

func (g GenGrpc) Process(importPath, protoPath string) error {
	if err := g.Setup(); err != nil {
		return err
	}

	if err := g.Gen(importPath, protoPath); err != nil {
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

	//if err := g.SetGoEnv(); err != nil {
	//	return err
	//}

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
	c := exec.Command("git", "checkout", "main")
	c.Dir = g.cache
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	c = exec.Command("git", "fetch", "origin", "main")
	c.Dir = g.cache
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	c = exec.Command("git", "reset", "origin/main", "--hard")
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
func (g GenGrpc) Gen(importPath, protoPath string) error {
	if err := os.MkdirAll("export", os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll("openapi", os.ModePerm); err != nil {
		return err
	}

	if err := g.protoC(importPath, protoPath); err != nil {
		return err
	}

	return nil
}

func (g GenGrpc) protoC(importPath, protoPath string) error {
	var pf []string
	var err error
	var prefix = protoPath
	if importPath != "" {
		tfp := filepath.Join(importPath, protoPath)
		if err := os.MkdirAll(tfp, os.ModePerm); err != nil {
			return err
		}
		defer func() {
			_ = os.RemoveAll(tfp)
		}()

		pf, err = CopyDir(protoPath, filepath.Join(importPath, protoPath))
		if err != nil {
			return err
		}
	} else {
		f, err := ioutil.ReadDir(protoPath)
		if err != nil {
			return err
		}
		for _, v := range f {
			pf = append(pf, v.Name())
		}
	}

	protoFile := make([]string, 0, len(pf))
	for _, v := range pf {
		protoFile = append(protoFile, filepath.Join(prefix, v))
	}

	args := make([]string, 0, 20)
	args = append(args, "protoc")

	if importPath != "" {
		args = append(args, "-I", importPath)
	}
	args = append(args,
		fmt.Sprintf("--go_out=./"),
		fmt.Sprintf("--go_opt=module=%s", g.module),
		fmt.Sprintf("--go-grpc_out=./"),
		fmt.Sprintf("--go-grpc_opt=module=%s", g.module),
		fmt.Sprintf("--go-gsv_out=./"),
		fmt.Sprintf("--go-gsv_opt=module=%s", g.module),
		fmt.Sprintf("--grpc-gateway_out=:./"),
		fmt.Sprintf("--grpc-gateway_opt=logtostderr=true"),
		fmt.Sprintf("--grpc-gateway_opt=module=%s", g.module),
		fmt.Sprintf("--openapiv2_out=./openapi"),
		fmt.Sprintf("--openapiv2_opt=logtostderr=true"),
		fmt.Sprintf("--openapiv2_opt=allow_merge=true"),
		fmt.Sprintf("--openapiv2_opt=merge_file_name=%s", filepath.Base(g.module)),
	)

	args = append(args, protoFile...)

	log.Println("[INFO] ", args)
	c := exec.Command(args[0], args[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return err
	}

	c = exec.Command("protoc-go-inject-tag", `-input=export/*/*.pb.go`)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	_ = c.Run()

	return nil
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
