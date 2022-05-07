package run

import (
	"github.com/kubetrail/ethkey/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type persistentFlagValues struct {
	OutputFormat string `json:"outputFormat,omitempty"`
}

func getPersistentFlags(cmd *cobra.Command) persistentFlagValues {
	rootCmd := cmd.Root().PersistentFlags()

	_ = viper.BindPFlag(flags.OutputFormat, rootCmd.Lookup(flags.OutputFormat))

	outputFormat := viper.GetString(flags.OutputFormat)

	return persistentFlagValues{
		OutputFormat: outputFormat,
	}
}
