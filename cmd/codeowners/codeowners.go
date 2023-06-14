package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.szostok.io/codeowners/internal/api"
	"go.szostok.io/codeowners/internal/config"
	"go.szostok.io/codeowners/internal/load"
	"go.szostok.io/codeowners/internal/runner"
	"go.szostok.io/codeowners/pkg/codeowners"
	"go.szostok.io/version/extension"
)

var severity api.SeverityType

// NewRoot returns a root cobra.Command for the whole Agent utility.
func RootCmd() *cobra.Command {
	cfg := &config.Config{}

	rootCmd := &cobra.Command{
		Use:          "codeowners",
		Short:        "Ensures the correctness of your CODEOWNERS file.",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return InitializeConfig(cmd, cfg, args)
		},
	}

	rootCmd.AddCommand(
		extension.NewVersionCobraCmd(),
		validateCmd(cfg),
	)

	return rootCmd
}

func validateCmd(cfg *config.Config) *cobra.Command {

	var validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate a CODEOWNERS file",

		Run: func(cmd *cobra.Command, args []string) {
			log := logrus.New()

			// init checks
			checks, err := load.Checks(cmd.Context(), cfg)
			exitOnError(err)

			// init codeowners entries
			codeownersEntries, err := codeowners.NewFromPath(cfg.RepositoryPath)
			exitOnError(err)

			// run check runner
			absRepoPath, err := filepath.Abs(cfg.RepositoryPath)
			exitOnError(err)

			checkRunner := runner.NewCheckRunner(log, codeownersEntries, absRepoPath, cfg.CheckFailureLevel, checks...)
			checkRunner.Run(cmd.Context())

			if cmd.Context().Err() != nil {
				log.Error("Application was interrupted by operating system")
				os.Exit(2)
			}
			if checkRunner.ShouldExitWithCheckFailure() {
				os.Exit(3)
			}
		},
	}
	addValidateFlags(validateCmd)
	return validateCmd
}

func addValidateFlags(cmd *cobra.Command) {
	cmd.Flags().StringSlice("checks", nil, "List of checks to be executed")
	cmd.Flags().Var(&severity, "check-failure-level", "Defines the level on which the application should treat check issues as failures")
	cmd.Flags().String("experimental-checks", "", "The comma-separated list of experimental checks that should be executed")
	cmd.Flags().String("github-access-token", "", "GitHub access token")
	cmd.Flags().String("github-base-url", "https://api.github.com/", "GitHub base URL for API requests")
	cmd.Flags().String("github-upload-url", "https://uploads.github.com/", "GitHub upload URL for uploading files")
	cmd.Flags().String("github-app-id", "", "Github App ID for authentication")
	cmd.Flags().String("github-app-installation-id", "", "Github App Installation ID")
	cmd.Flags().String("github-app-private-key", "", "Github App private key in PEM format")
	cmd.Flags().StringSlice("not-owned-checker-skip-patterns", nil, "The comma-separated list of patterns that should be ignored by not-owned-checker")
	cmd.Flags().StringSlice("not-owned-checker-subdirectories", nil, "The comma-separated list of subdirectories to check in not-owned-checker")
	cmd.Flags().Bool("not-owned-checker-trust-workspace", false, "Specifies whether the repository path should be marked as safe")
	cmd.Flags().String("repository-path", "", "Path to your repository on your local machine")
	cmd.Flags().String("owner-checker-repository", "", "The owner and repository name separated by slash")
	cmd.Flags().StringSlice("owner-checker-ignored-owners", []string{"@ghost"}, "The comma-separated list of owners that should not be validated")
	cmd.Flags().Bool("owner-checker-allow-unowned-patterns", true, "Specifies whether CODEOWNERS may have unowned files")
	cmd.Flags().Bool("owner-checker-owners-must-be-teams", false, "Specifies whether only teams are allowed as owners of files")
}

func exitOnError(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

func InitializeConfig(cmd *cobra.Command, cfg *config.Config, args []string) error {
	v := viper.New()

	// Look for config file, ignore if missing
	v.SetConfigName(config.DefaultConfigFilename)
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	// Look for environment variables
	v.SetEnvPrefix(config.EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Bind flags to the configuration struct
	bindFlags(cmd, v)

	// Unmarshal the configuration into the struct
	if err := v.Unmarshal(cfg); err != nil {
		return err
	}

	return nil
}

// Bind each cobra flag to its associated viper configuration environment variable
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
		v.BindPFlag(configName, f)
	})
}
