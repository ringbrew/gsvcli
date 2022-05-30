package subcmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CopyDir(src, dst string, filter ...func(name string) bool) ([]string, error) {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return nil, err
	}

	copied := make([]string, 0, len(files))

	for _, v := range files {
		if !v.IsDir() {
			if len(filter) > 0 && filter[0] != nil {
				if !filter[0](v.Name()) {
					continue
				}
			}

			if err := CopyFile(filepath.Join(src, v.Name()), filepath.Join(dst, v.Name())); err != nil {
				return nil, err
			}

			copied = append(copied, v.Name())
		}
	}

	return copied, nil
}

func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func GoFmt(dir ...string) error {
	c := exec.Command("go", "fmt", "./...")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if len(dir) > 0 && dir[0] != "" {
		c.Dir = dir[0]
	}
	if err := c.Run(); err != nil {
		return err
	}
	return nil
}
