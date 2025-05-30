package cmdenv

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/environments/envproviders"
	"slv.sh/slv/internal/core/input"
	"slv.sh/slv/internal/core/profiles"
)

func envNewCommand() *cobra.Command {
	if envNewCmd == nil {
		envNewCmd = &cobra.Command{
			Use:   "new",
			Short: "Create a new environment",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Help()
			},
		}
		envNewCmd.PersistentFlags().BoolP(utils.QuantumSafeFlag.Name, utils.QuantumSafeFlag.Shorthand, false, utils.QuantumSafeFlag.Usage)
		envNewCmd.AddCommand(envNewServiceCommand())
		envNewCmd.AddCommand(envNewUserCommand())
	}
	return envNewCmd
}

func envNewServiceCommand() *cobra.Command {
	if envNewServiceCmd == nil {
		envNewServiceCmd = &cobra.Command{
			Use:   "service",
			Short: "Creates a new service environment",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Help()
			},
		}
		envNewServiceCmd.PersistentFlags().StringP(envNameFlag.Name, envNameFlag.Shorthand, "", envNameFlag.Usage)
		envNewServiceCmd.PersistentFlags().StringP(envEmailFlag.Name, envEmailFlag.Shorthand, "", envEmailFlag.Usage)
		envNewServiceCmd.PersistentFlags().StringSliceP(envTagsFlag.Name, envTagsFlag.Shorthand, []string{}, envTagsFlag.Usage)
		envNewServiceCmd.PersistentFlags().BoolP(envAddFlag.Name, envAddFlag.Shorthand, false, envAddFlag.Usage)
		envNewServiceCmd.MarkPersistentFlagRequired(envNameFlag.Name)
		envNewServiceCmd.AddCommand(envNewServicePlaintextCommand())
		envNewServiceCmd.AddCommand(newKMSEnvCommand("aws", "Create a service environment for AWS KMS", awsARNFlag))
		envNewServiceCmd.AddCommand(newKMSEnvCommand("gcp", "Create a service environment for GCP KMS", gcpKmsResNameFlag))
	}
	return envNewServiceCmd
}

func envNewServicePlaintextCommand() *cobra.Command {
	if envNewServicePlaintextCmd == nil {
		envNewServicePlaintextCmd = &cobra.Command{
			Use:     "plaintext",
			Aliases: []string{"direct", "raw"},
			Short:   "Creates a new service environment and returns the secret key in plaintext",
			Run: func(cmd *cobra.Command, args []string) {
				name, _ := cmd.Flags().GetString(envNameFlag.Name)
				email, _ := cmd.Flags().GetString(envEmailFlag.Name)
				tags, err := cmd.Flags().GetStringSlice(envTagsFlag.Name)
				if err != nil {
					utils.ExitOnError(err)
				}
				addToProfileFlag, _ := cmd.Flags().GetBool(envAddFlag.Name)
				var profile *profiles.Profile
				if addToProfileFlag {
					profile, err = profiles.GetActiveProfile()
					if err != nil {
						utils.ExitOnError(err)
					}
					if !profile.IsPushSupported() {
						utils.ExitOnError(fmt.Errorf("profile (%s) does not support adding environments", profile.Name()))
					}
				}
				var env *environments.Environment
				var secretKey *crypto.SecretKey
				pq, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
				env, secretKey, err = environments.NewEnvironment(name, environments.SERVICE, pq)
				if err != nil {
					utils.ExitOnError(err)
				}
				env.SetEmail(email)
				env.AddTags(tags...)
				ShowEnv(*env, true, false)
				if secretKey != nil {
					fmt.Println("\nSecret Key:\t", color.HiBlackString(secretKey.String()))
				}
				if addToProfileFlag {
					if err = profile.PutEnv(env); err != nil {
						utils.ExitOnError(fmt.Errorf("failed to add the environment to profile (%s): %w", profile.Name(), err))
					}
					fmt.Printf("Successfully added the environment to profile (%s)\n", color.GreenString(profile.Name()))
				}
			},
		}
	}
	return envNewServicePlaintextCmd
}

func envNewUserCommand() *cobra.Command {
	if envNewUserCmd == nil {
		envNewUserCmd = &cobra.Command{
			Use:     "self",
			Aliases: []string{"user", "usr", "u"},
			Short:   "Register as a new user environment",
			Run: func(cmd *cobra.Command, args []string) {
				selfEnv := environments.GetSelf()
				if selfEnv != nil {
					ShowEnv(*selfEnv, true, true)
					confirmed, err := input.GetConfirmation("You are already registered as an environment, "+
						"this will replace the existing environment. Proceed? (yes/no): ", "yes")
					if err != nil {
						utils.ExitOnError(err)
					}
					if !confirmed {
						fmt.Println(color.YellowString("Operation aborted"))
						utils.SafeExit()
					}
				}
				addToProfileFlag, _ := cmd.Flags().GetBool(envAddFlag.Name)
				var err error
				var profile *profiles.Profile
				if addToProfileFlag {
					profile, err = profiles.GetActiveProfile()
					if err != nil {
						utils.ExitOnError(err)
					}
					if !profile.IsPushSupported() {
						utils.ExitOnError(fmt.Errorf("profile (%s) does not support adding environments", profile.Name()))
					}
				}
				envName, _ := cmd.Flags().GetString(envNameFlag.Name)
				envEmail, _ := cmd.Flags().GetString(envEmailFlag.Name)
				envTags, err := cmd.Flags().GetStringSlice(envTagsFlag.Name)
				if err != nil {
					utils.ExitOnError(err)
				}
				inputs := make(map[string][]byte)
				password, err := input.NewPasswordFromUser(input.DefaultPasswordPolicy())
				if err != nil {
					utils.ExitOnError(err)
				}
				inputs["password"] = password
				var env *environments.Environment
				pq, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
				env, err = envproviders.NewEnv("password", envName, environments.USER, inputs, pq)
				if err != nil {
					utils.ExitOnError(err)
				}
				env.SetEmail(envEmail)
				env.AddTags(envTags...)
				if err = env.MarkAsSelf(); err != nil {
					utils.ExitOnError(err)
				}
				secretBinding := env.SecretBinding
				ShowEnv(*env, true, true)
				if addToProfileFlag {
					if err = profile.PutEnv(env); err != nil {
						utils.ExitOnError(fmt.Errorf("failed to add the environment to profile (%s): %w", profile.Name(), err))
					}
					fmt.Printf("Successfully added the environment to profile (%s)\n", color.GreenString(profile.Name()))
				}
				fmt.Println(color.GreenString("Successfully registered as self environment"))
				if secretBinding != "" {
					fmt.Println(color.YellowString("Please note down the \"Secret Binding\" somewhere safe so that you don't lose it.\n" +
						"It is required to access your registered environment."))
				}
				utils.SafeExit()
			},
		}
		envNewUserCmd.Flags().StringP(envNameFlag.Name, envNameFlag.Shorthand, "", envNameFlag.Usage)
		envNewUserCmd.Flags().StringP(envEmailFlag.Name, envEmailFlag.Shorthand, "", envEmailFlag.Usage)
		envNewUserCmd.Flags().StringSliceP(envTagsFlag.Name, envTagsFlag.Shorthand, []string{}, envTagsFlag.Usage)
		envNewUserCmd.Flags().BoolP(envAddFlag.Name, envAddFlag.Shorthand, false, envAddFlag.Usage)
		envNewUserCmd.MarkFlagRequired(envNameFlag.Name)
	}
	return envNewUserCmd
}
