package cmd

import (
	"fmt"

	"github.com/meowfaceman/conshim/pkg/shims"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists current shims.",
		Long:  "Lists all of the current shims and displays their associated container commands.",

		Run: func(cmd *cobra.Command, args []string) {
			shims, err := shims.List()
			cobra.CheckErr(err)

			fmt.Printf("%25s    %7s   %s\n", "Name", "Source", "Version")
			fmt.Println("--------------------------|---------|---------")
			for _, shim := range shims {
				fmt.Printf("%25s    %7s   %s\n", shim.Name, shim.Source, shim.Version)
			}
		},
	}
)
