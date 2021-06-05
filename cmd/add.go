package cmd

import (
	"fmt"
	"strings"

	"github.com/meowfaceman/conshim/pkg/shims"
	"github.com/spf13/cobra"
)

var (
	addShimName    string
	addShimCommand string

	addCmd = &cobra.Command{
		Use:   "add <shim> <command>",
		Short: "Adds a shim.",
		Long:  "Adds a shim and attaches the associated command with it.",

		Args: func(cmd *cobra.Command, args []string) error {
			numArgs := len(args)
			if numArgs < 2 {
				return fmt.Errorf("required 2 or more arguments, got %d", numArgs)
			}

			addShimName = args[0]
			addShimCommand = strings.Join(args[1:], " ")

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			err := shims.Add(addShimName, addShimCommand)
			cobra.CheckErr(err)
		},
	}
)
