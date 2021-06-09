package cmd

import (
	"fmt"

	"github.com/meowfaceman/conshim/pkg/config"
	"github.com/spf13/cobra"
)

var (
	binPathCmd = &cobra.Command{
		Use:   "binpath",
		Short: "Prints out the bin path.",
		Long:  "Prints out the bin path where the shims are located.",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.BinPath())
		},
	}
)
