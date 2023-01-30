package main

import (
	"github.com/morlay/crawler/pkg/crawler/server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use: "serve",
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Serve(cmd.Context(), crawlerInst)
	},
}
