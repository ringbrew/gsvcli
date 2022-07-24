package subcmd

import (
	"errors"
	"github.com/spf13/cobra"
	"log"
)

type GenDomainSubCmd struct {
	Domain string
}

func (sc *GenDomainSubCmd) Parse(args []string) error {
	if len(args) == 0 {
		return errors.New("invalid domain")
	}
	sc.Domain = args[0]
	return nil
}

func (sc *GenDomainSubCmd) Process() error {
	log.Printf("[INFO] gen domain command is running\n")

	p := NewGenDomain(sc.Domain)
	return p.Process()
}

func NewDomainCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "domain",
		Short: "generator for domain code",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal(errors.New("invalid domain").Error())
			}
			domain := args[0]

			log.Printf("gen domain command is running domain[%s]\n", domain)

			p := NewGenDomain(domain)
			if err := p.Process(); err != nil {
				log.Fatal(err.Error())
			}
		},
	}
}
