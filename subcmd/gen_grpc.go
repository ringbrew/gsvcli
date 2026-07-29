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
	protoFile, err := listProtoFile(protoPath)
	if err != nil {
		return err
	}

	if len(protoFile) == 0 {
		return fmt.Errorf("no .proto file found under %s", protoPath)
	}

	args := make([]string, 0, 20)
	args = append(args, "protoc")

	// importPath(未显式指定时取环境变量 GOPROTO)是第三方依赖 proto 的集中存放目录，
	// 例如 google/api/annotations.proto、protoc-gen-openapiv2/options/annotations.proto。
	if importPath != "" {
		args = append(args, "-I", importPath)
	}

	// 项目自身的 proto 以仓库根目录为 include 根，descriptor 名因此保留 protoPath
	// 前缀(如 proto/account.proto)。这个名字是 protoregistry 的 key，也是 import
	// 语句必须写的路径，跨版本必须稳定 —— 用 -I protoPath 会把前缀剥掉。
	args = append(args, "-I", ".")

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
		fmt.Sprintf("--openapiv2_opt=enums_as_ints=true"),
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

	injectTag()

	return nil
}

// listProtoFile 递归收集 protoPath 下的 .proto，返回相对当前目录的路径。
// 递归是为了支持 proto/<domain>/<version>/x.proto 这类分层布局；过滤扩展名是
// 为了不把目录项和 README 之类的文件当成 proto 传给 protoc。
func listProtoFile(protoPath string) ([]string, error) {
	result := make([]string, 0, 8)

	if err := filepath.Walk(protoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".proto" {
			return nil
		}

		result = append(result, path)
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

// injectTag 把 proto 注释里的 // @gotags: 回填成 struct tag。protoc-go-inject-tag
// 的 -input 走 filepath.Glob，不支持 **，所以按 pb.go 所在目录逐个调用，
// 才能覆盖 export 下的多级目录。失败不阻断生成，与历史行为保持一致。
func injectTag() {
	dir := make(map[string]struct{})

	_ = filepath.Walk("export", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".pb.go") {
			return nil
		}

		dir[filepath.Dir(path)] = struct{}{}
		return nil
	})

	for d := range dir {
		c := exec.Command("protoc-go-inject-tag", fmt.Sprintf("-input=%s/*.pb.go", d))
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		_ = c.Run()
	}
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
