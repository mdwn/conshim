package shim

import (
	"fmt"

	"github.com/meowfaceman/conshim/pkg/config"
	"github.com/meowfaceman/conshim/pkg/shim"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists current shims.",
		Long:  "Lists all of the current shims and displays their associated container commands.",

		Run: func(cmd *cobra.Command, args []string) {
			shims, err := config.ListShims()
			cobra.CheckErr(err)

			fmt.Print(shim.ShimsListToString(shims))
		},
	}
)
