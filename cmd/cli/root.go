package cli

import (
	"fmt"
	"os"

	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/manifest-cyber/ai-bom/cmd/cli/options"
	"github.com/manifest-cyber/ai-bom/pkg/log"
)

var (
	version = "0.0.0-dev"
	name    = "ai-bom"
)

func ManCommand(root *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "man",
		Short:                 "Generates command line manpages",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			manPage, err := mcobra.NewManPage(1, root)
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
			return err
		},
	}

	return cmd
}

func NewRootCmd() *cobra.Command {
	ro := &options.RootOptions{}
	rootCmd := &cobra.Command{
		Use:     "ai-bom",
		Version: version,
		Short:   "AI-BOM generates BOMs for AI/ML models",
		Long:    "AI-BOM generates BOMs for AI/ML models",
		Run:     func(cmd *cobra.Command, args []string) {},
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := validateRootOptions(ro); err != nil {
				return err
			}

			if err := options.BindConfig(viper.GetViper(), cmd); err != nil {
				return err
			}
			return setupLogger(ro)
		},
		SilenceErrors: true,
	}

	rootCmd.SetVersionTemplate(fmt.Sprintf("%s v{{.Version}}\n", name))

	ro.AddFlags(rootCmd)

	// Commands
	cvtCmd := BomCommand()
	rootCmd.AddCommand(cvtCmd)

	// Manpages
	rootCmd.AddCommand(ManCommand(rootCmd))
	return rootCmd
}

func validateRootOptions(_ *options.RootOptions) error {
	return nil
}

func setupLogger(ro *options.RootOptions) error {
	level := zapcore.Level(int(zap.WarnLevel) - ro.Verbose)
	log, err := log.NewLogger(
		log.WithLevel(level),
		log.WithGlobalLogger(),
	)
	if err != nil {
		return err
	}

	log.Debug("logger initialized")
	return nil
}
