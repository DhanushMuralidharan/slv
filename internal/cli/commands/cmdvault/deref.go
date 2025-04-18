package cmdvault

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/secretkey"
)

func vaultDerefCommand() *cobra.Command {
	if vaultDerefCmd == nil {
		vaultDerefCmd = &cobra.Command{
			Use:   "deref",
			Short: "Dereferences and updates secrets from a vault to a given yaml or json file",
			Run: func(cmd *cobra.Command, args []string) {
				envSecretKey, err := secretkey.Get()
				if err != nil {
					utils.ExitOnError(err)
				}
				vaultFiles, err := cmd.Flags().GetStringSlice(vaultFileFlag.Name)
				if err != nil {
					utils.ExitOnError(err)
				}
				paths, err := cmd.Flags().GetStringSlice(vaultDerefPathFlag.Name)
				if err != nil {
					utils.ExitOnError(err)
				}
				for _, vaultFile := range vaultFiles {
					vault, err := getVault(vaultFile)
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
		vaultDerefCmd.Flags().StringSliceP(vaultDerefPathFlag.Name, vaultDerefPathFlag.Shorthand, []string{}, vaultDerefPathFlag.Usage)
		vaultDerefCmd.MarkFlagRequired(vaultDerefPathFlag.Name)
	}
	return vaultDerefCmd
}
