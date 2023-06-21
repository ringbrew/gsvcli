package main

import (
	"github.com/ringbrew/gsvcli/subcmd"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gsv",
		Short: "gsv helper command tools.",
	}

	genCmd := &cobra.Command{
		Use:   "gen",
		Short: "gsv code generator.",
	}

	rootCmd.AddCommand(subcmd.NewInitCommand())
	rootCmd.AddCommand(subcmd.NewInstallCommand())

	genCmd.AddCommand(subcmd.NewGrpcCommand())
	genCmd.AddCommand(subcmd.NewDomainCommand())
	genCmd.AddCommand(subcmd.NewHttpCommand())
	genCmd.AddCommand(subcmd.NewHandlerCommand())
	genCmd.AddCommand(subcmd.NewGenProtoSubCmd())

	rootCmd.AddCommand(genCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
