package cmdvault

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/session"
	"slv.sh/slv/internal/core/vaults"
)

func vaultDerefCommand() *cobra.Command {
	if vaultDerefCmd == nil {
		vaultDerefCmd = &cobra.Command{
			Use:   "deref",
			Short: "Dereferences and updates values from a vault to a given file with vault references",
			Run: func(cmd *cobra.Command, args []string) {
				envSecretKey, err := session.GetSecretKey()
				if err != nil {
					utils.ExitOnError(err)
				}
				vaultFiles, err := cmd.Flags().GetStringSlice(vaultFileFlag.Name)
				if err != nil {
					utils.ExitOnError(err)
				}
				paths, err := cmd.Flags().GetStringSlice(vaultRefFileFlag.Name)
				if err != nil {
					utils.ExitOnError(err)
				}
				for _, vaultFile := range vaultFiles {
					vault, err := vaults.Get(vaultFile)
					if err != nil {
						utils.ExitOnError(err)
					}
					err = vault.Unlock(envSecretKey)
					if err != nil {
						utils.ExitOnError(err)
					}
					for _, path := range paths {
						if err = vault.DeRef(path); err != nil {
							utils.ExitOnError(err)
						}
						fmt.Println("Dereferenced", color.GreenString(path), "with the vault", color.GreenString(vaultFile))
					}
				}
				utils.SafeExit()
			},
		}
		vaultDerefCmd.Flags().StringSliceP(vaultFileFlag.Name, vaultFileFlag.Shorthand, []string{}, vaultFileFlag.Usage)
		vaultDerefCmd.Flags().StringSliceP(vaultRefFileFlag.Name, vaultRefFileFlag.Shorthand, []string{}, vaultRefFileFlag.Usage)
		vaultDerefCmd.MarkFlagRequired(vaultRefFileFlag.Name)
	}
	return vaultDerefCmd
}
