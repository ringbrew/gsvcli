package subcmd

import (
	"errors"
	"github.com/spf13/cobra"
	"log"
)

type GenHandlerSubCmd struct {
	Domain string
}

func (sc *GenHandlerSubCmd) Parse(args []string) error {
	if len(args) == 0 {
		return errors.New("invalid domain")
	}
	sc.Domain = args[0]
	return nil
}

func (sc *GenHandlerSubCmd) Process() error {
	log.Printf("[INFO] gen handler command is running\n")

	p := NewGenHttp(sc.Domain)
	return p.Process()
}

func NewHandlerCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "handler",
		Short:   "generator for handler code. usage: gsv gen handler {domain} optional:{subdomain}",
		Example: "gsv gen http demo example",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 && len(args) != 2 {
				log.Fatal(errors.New("invalid input").Error())
			}
			domain := args[0]
			subdomain := domain
			if len(args) == 2 {
				subdomain = args[1]
			}

			log.Printf("gen handler command is running domain[%s]\n", domain)

			p := NewGenHandler(domain, subdomain)
			if err := p.Process(); err != nil {
				log.Fatal(err.Error())
			}
		},
	}
}
