package subcmd

import "errors"

type SubCmd interface {
	Parse(args []string) error
	Process() error
}

func New(name Name) (SubCmd, error) {
	genCmd, exist := subCmdManager[name]
	if !exist {
		return nil, errors.New("not support sub command")
	}

	return genCmd(), nil
}

type Name string

const (
	Init Name = "init"
	Gen  Name = "gen"
)

var subCmdManager = map[Name]func() SubCmd{}

func Register(name Name, cmdGen func() SubCmd) {
	subCmdManager[name] = cmdGen
}
