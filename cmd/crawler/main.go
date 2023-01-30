package main

import (
	"context"
	"log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/spf13/cobra"

	"github.com/morlay/crawler/pkg/crawler"
)

var spec string
var verbose int
var crawlerInst crawler.Crawler

var rootCmd = &cobra.Command{
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		spec, err := os.ReadFile(spec)
		if err != nil {
			return err
		}

		s, err := crawler.NewCrawler(spec)
		if err != nil {
			return err
		}
		crawlerInst = s
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().
		StringVarP(&spec, "spec", "s", "source.gql", "source spec")
	rootCmd.PersistentFlags().
		IntVarP(&verbose, "verbose", "v", int(stdr.Info), "log verbose")
}

func main() {
	stdr.SetVerbosity(verbose)

	ctx := logr.NewContext(context.Background(), stdr.New(log.New(os.Stdout, "[crawler] ", log.Ldate|log.Ltime)))

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
