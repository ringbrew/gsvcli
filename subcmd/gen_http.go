package subcmd

import (
	"bufio"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type GenHttp struct {
	module string
	domain string
}

const httpServiceGenImpl = `package [[.domain]]

import (
	"github.com/ringbrew/gsv/service"
	"[[.module]]/internal/domain"
)

type Service struct {
	ctx *domain.UseCaseContext

	name   string
	remark string
	desc   service.Description
}

func NewService(ctx *domain.UseCaseContext) service.Service {
	return &Service{
		ctx: ctx,
		name:   "api.[[.domain]].service",
		remark: "",
	}
}

func (s *Service) Name() string {
	return s.name
}

func (s *Service) Remark() string {
	return s.remark
}

func (s *Service) Description() service.Description {
	return s.desc
}`

func NewGenHttp(domain string) *GenHttp {
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
		log.Fatal("[FATAL] " + "please init your go module with go mod init")
	}

	return &GenHttp{
		module: module,
		domain: domain,
	}
}

func (gh *GenHttp) Process() error {
	deliveryPath := filepath.Join("delivery", gh.domain, "service.http.impl.go")

	if _, err := os.Stat(deliveryPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	} else {
		return nil
	}

	if err := os.MkdirAll(deliveryPath, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(deliveryPath)
	if err != nil {
		return err
	}

	t, err := template.New("service.http.impl.go").Parse(httpServiceGenImpl)
	if err != nil {
		return err
	}

	if err := t.Execute(file, map[string]interface{}{
		"module": gh.module,
		"domain": gh.domain,
	}); err != nil {
		return err
	}

	c := exec.Command("go", "fmt")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Dir = deliveryPath
	if err := c.Run(); err != nil {
		return err
	}

	return nil
}
