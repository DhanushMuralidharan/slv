package cmdvault

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/cmdenv"
	"oss.amagi.com/slv/internal/cli/commands/utils"
)

func vaultNewCommand() *cobra.Command {
	if vaultNewCmd == nil {
		vaultNewCmd = &cobra.Command{
			Use:   "new",
			Short: "Creates a new vault",
			Run: func(cmd *cobra.Command, args []string) {
				vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
				pq, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
				publicKeys, err := cmdenv.GetPublicKeys(cmd, true, pq)
				if err != nil {
					utils.ExitOnError(err)
				}
				enableHash, _ := cmd.Flags().GetBool(vaultEnableHashingFlag.Name)
				k8sName := cmd.Flag(vaultK8sNameFlag.Name).Value.String()
				k8sNamespace := cmd.Flag(vaultK8sNamespaceFlag.Name).Value.String()
				k8sSecret := cmd.Flag(vaultK8sSecretFlag.Name).Value.String()
				if _, err = newK8sVault(vaultFile, k8sName, k8sNamespace, k8sSecret, enableHash, pq, publicKeys...); err != nil {
					utils.ExitOnError(err)
				}
				fmt.Println("Created vault:", color.GreenString(vaultFile))
				utils.SafeExit()
			},
		}
		vaultNewCmd.Flags().StringSliceP(cmdenv.EnvPublicKeysFlag.Name, cmdenv.EnvPublicKeysFlag.Shorthand, []string{}, cmdenv.EnvPublicKeysFlag.Usage)
		vaultNewCmd.Flags().StringSliceP(cmdenv.EnvSearchFlag.Name, cmdenv.EnvSearchFlag.Shorthand, []string{}, cmdenv.EnvSearchFlag.Usage)
		vaultNewCmd.Flags().BoolP(cmdenv.EnvSelfFlag.Name, cmdenv.EnvSelfFlag.Shorthand, false, cmdenv.EnvSelfFlag.Usage)
		vaultNewCmd.Flags().BoolP(cmdenv.EnvK8sFlag.Name, cmdenv.EnvK8sFlag.Shorthand, false, cmdenv.EnvK8sFlag.Usage)
		vaultNewCmd.Flags().StringP(vaultK8sNameFlag.Name, vaultK8sNameFlag.Shorthand, "", vaultK8sNameFlag.Usage)
		vaultNewCmd.Flags().StringP(vaultK8sNamespaceFlag.Name, vaultK8sNamespaceFlag.Shorthand, "", vaultK8sNamespaceFlag.Usage)
		vaultNewCmd.Flags().StringP(vaultK8sSecretFlag.Name, vaultK8sSecretFlag.Shorthand, "", vaultK8sSecretFlag.Usage)
		vaultNewCmd.Flags().BoolP(vaultEnableHashingFlag.Name, vaultEnableHashingFlag.Shorthand, false, vaultEnableHashingFlag.Usage)
		vaultNewCmd.Flags().BoolP(utils.QuantumSafeFlag.Name, utils.QuantumSafeFlag.Shorthand, false, utils.QuantumSafeFlag.Usage)
	}
	return vaultNewCmd
}
