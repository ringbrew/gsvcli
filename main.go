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
		Use: "gen",
	}

	rootCmd.AddCommand(subcmd.NewInitCommand())
	rootCmd.AddCommand(subcmd.NewInstallCommand())

	genCmd.AddCommand(subcmd.NewGrpcCommand())
	genCmd.AddCommand(subcmd.NewDomainCommand())
	rootCmd.AddCommand(genCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
