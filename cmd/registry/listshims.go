package registry

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	listShimsCmdRegistryName string

	listShimsCmd = &cobra.Command{
		Use:   "list-shims <registry>",
		Short: "Lists the shims from the given registry.",
		Long:  "Lists the shims from the given registry.",

		Args: func(cmd *cobra.Command, args []string) error {
			numArgs := len(args)

			if len(args) != 1 {
				return fmt.Errorf("expected 1 arguments, got %d", numArgs)
			}

			listShimsCmdRegistryName = args[0]

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			m, err := addOrGetRegistry(listShimsCmdRegistryName)
			cobra.CheckErr(err)

			fmt.Print(m.ShimsToString())
		},
	}
)
