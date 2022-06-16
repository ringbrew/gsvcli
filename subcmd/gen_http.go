package subcmd

import (
	"bufio"
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
	tmpl   map[string]string
}

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
		tmpl: map[string]string{
			"service.http.impl.go": httpGenTmpl,
		},
	}
}

func (gd *GenHttp) Process() error {
	domainPath := filepath.Join("internal", "delivery", gd.domain)

	if err := os.MkdirAll(domainPath, os.ModePerm); err != nil {
		return err
	}

	fileList := make([]*os.File, 0, len(gd.tmpl))
	defer func() {
		for i := range fileList {
			fileList[i].Close()
		}
	}()

	for fileName, tmpl := range gd.tmpl {
		fileFullPath := filepath.Join(domainPath, fileName)

		if _, err := os.Stat(fileFullPath); err != nil && os.IsNotExist(err) {
			file, err := os.Create(fileFullPath)
			if err != nil {
				return err
			}
			fileList = append(fileList, file)

			t, err := template.New(fileName).Delims("[[", "]]").Parse(tmpl)
			if err != nil {
				return err
			}

			if err := t.Execute(file, map[string]interface{}{
				"module":      gd.module,
				"serviceName": "Service",
				"packageName": gd.domain,
			}); err != nil {
				return err
			}
		}
	}

	c := exec.Command("go", "fmt")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Dir = domainPath
	if err := c.Run(); err != nil {
		return err
	}

	return nil
}

const httpGenTmpl = `package [[.packageName]]

import (
	"github.com/ringbrew/gsv/service"
	"[[.module]]/internal/domain"
)

type [[.serviceName]] struct {
	ctx *domain.UseCaseContext
	name   string
	remark string
	desc   service.Description
	service.HttpRouteCollector
}

func New[[.serviceName]](ctx *domain.UseCaseContext) service.Service {
	s := &[[.serviceName]]{
		ctx:    ctx,
		name:   "api.[[.packageName]].service",
		remark: "",
		desc: service.Description{
			Valid: true,
		},
	}
	s.desc.HttpRoute = s.HttpRouteCollector
	return s
}

func (s *[[.serviceName]]) Name() string {
	return s.name
}

func (s *[[.serviceName]]) Remark() string {
	return s.remark
}

func (s *[[.serviceName]]) Description() service.Description {
	return s.desc
}`
