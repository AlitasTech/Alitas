package commands

import (
	"github.com/spf13/cobra"
)

var (
	_config = NewDefaultCLIConf()
)

// RootCmd is the root command for Alitas Network
var RootCmd = &cobra.Command{
	Use:              "Alitas Network",
	Short:            "alitas",
	TraverseChildren: true,
}
