package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fornellas/tasmota_exporter/version"
)

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the program version.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", version.GetVersion())
	},
}
