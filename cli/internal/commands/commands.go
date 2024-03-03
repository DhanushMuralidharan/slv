package commands

import "github.com/spf13/cobra"

var (

	// SLV Command
	slvCmd *cobra.Command

	// Version Command
	versionCmd *cobra.Command

	// System Commands
	systemCmd      *cobra.Command
	systemResetCmd *cobra.Command

	// Profile Commands
	profileCmd     *cobra.Command
	profileNewCmd  *cobra.Command
	profileListCmd *cobra.Command
	profileSetCmd  *cobra.Command
	profileDelCmd  *cobra.Command
	profilePullCmd *cobra.Command
	profilePushCmd *cobra.Command

	// Environment Commands
	envCmd           *cobra.Command
	envNewCmd        *cobra.Command
	envNewServiceCmd *cobra.Command
	envNewUserCmd    *cobra.Command
	envAddCmd        *cobra.Command
	envListSearchCmd *cobra.Command
	envSelfCmd       *cobra.Command
	envSelfSetCmd    *cobra.Command

	// Vault Commands
	vaultCmd       *cobra.Command
	vaultNewCmd    *cobra.Command
	vaultShareCmd  *cobra.Command
	vaultInfoCmd   *cobra.Command
	vaultPutCmd    *cobra.Command
	vaultGetCmd    *cobra.Command
	vaultExportCmd *cobra.Command
	vaultRefCmd    *cobra.Command
	vaultDerefCmd  *cobra.Command

	// Secret Commands
	// secretCmd *cobra.Command
)
