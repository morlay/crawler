package main

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var parameters []string

func init() {
	doCmd.Flags().StringSliceVarP(&parameters, "parameter", "p", nil, "parameter <name>=<value>")

	rootCmd.AddCommand(doCmd)
}

var doCmd = &cobra.Command{
	Use: "do [operation]",
	RunE: func(cmd *cobra.Command, args []string) error {
		o := args[0]
		params := map[string]string{}

		for _, p := range parameters {
			parts := strings.SplitN(p, "=", 2)
			if len(parts) == 2 {
				params[parts[0]] = parts[1]
			} else {
				params[parts[0]] = ""
			}
		}

		ctx := cmd.Context()

		ret, err := crawlerInst.Do(ctx, o, params)
		if err != nil {
			return err
		}

		return ret.Scan(ctx, os.Stdout)
	},
}
