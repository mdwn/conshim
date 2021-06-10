package shim

import (
	"fmt"
	"strings"

	"github.com/meowfaceman/conshim/pkg/config"
	"github.com/spf13/cobra"
)

var (
	updateShimName    string
	updateShimCommand string

	updateCmd = &cobra.Command{
		Use:   "update <shim> <command>",
		Short: "Updates a shim.",
		Long:  "Updates an existing shim and attaches the associated command with it.",

		Args: func(cmd *cobra.Command, args []string) error {
			numArgs := len(args)
			if numArgs < 2 {
				return fmt.Errorf("required 2 or more arguments, got %d", numArgs)
			}

			updateShimName = args[0]
			updateShimCommand = strings.Join(args[1:], " ")

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			err := config.UpdateShim(updateShimName, "user", "NONE", []string{}, updateShimCommand)
			cobra.CheckErr(err)
		},
	}
)
