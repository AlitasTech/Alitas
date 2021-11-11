package commands

import (
	"fmt"

	"github.com/AlitasTech/Alitas/src/version"
	"github.com/spf13/cobra"
)

// VersionCmd displays the version of alitas being used
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show alitas version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version)
	},
}
