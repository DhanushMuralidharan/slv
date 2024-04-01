package cmdsystem

import (
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/cli/internal/commands/utils"
)

var (
	systemCmd      *cobra.Command
	systemResetCmd *cobra.Command
)

var (
	// Common Flags
	yesFlag = utils.FlagDef{
		Name:      "yes",
		Shorthand: "y",
		Usage:     "Confirm action",
	}
)
