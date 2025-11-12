package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chigaji/diag-cli/internal/collectors/gopsutil"
	"github.com/chigaji/diag-cli/internal/config"
	"github.com/chigaji/diag-cli/internal/log"
	"github.com/chigaji/diag-cli/internal/render"
	"github.com/chigaji/diag-cli/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	output  string
	noColor bool
	logLvl  string
)

func main() {

	root := &cobra.Command{
		Use:   "diag",
		Short: "A system diagnostics CLI tool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// load config
			if err := config.Load(cfgFile); err != nil {
				return err
			}
			// flag overrides
			if output != "" {
				viper.Set("output", output)
			}

			if logLvl != "" {
				viper.Set("log.level", logLvl)
			}
			viper.Set("ui.no_color", noColor)

			// logger
			log.SetUp(viper.GetString("log.level"), viper.GetBool("ui.no_color"))
			return nil
		},
	}

	root.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path (default is $HOME/.diag.yaml)")
	root.PersistentFlags().StringVarP(&output, "output", "o", "", "output format: table|json|yaml")
	root.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable color output")
	root.PersistentFlags().StringVar(&logLvl, "log-level", "", "log level: debug|info|warn|error")

	// root.AddCommand()
	root.AddCommand(newSysCmd())
	root.AddCommand(newNetCmd())
	root.AddCommand(newProcCmd())
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.FullVersion())
		},
	})

	if err := root.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func newSysCmd() *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "sys",
		Short: "Show system metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := gopsutil.New()
			ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
			defer cancel()

			data, err := c.System(ctx, all)
			if err != nil {
				return err
			}
			return render.Print(cmd.OutOrStdout(), data)

		},
	}
	cmd.Flags().BoolVar(&all, "all", false, "collect all system metrics (temps, partitions, etc.)")
	return cmd
}

func newNetCmd() *cobra.Command {
	var iface string
	cmd := &cobra.Command{
		Use:   "net",
		Short: "Show network diagnostics",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := gopsutil.New()
			ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
			defer cancel()

			data, err := c.Network(ctx, iface)
			if err != nil {
				return err
			}
			return render.Print(cmd.OutOrStdout(), data)

		},
	}
	cmd.Flags().StringVar(&iface, "interface", "", "limit to interface name")
	return cmd
}

func newProcCmd() *cobra.Command {
	var top int
	var sortKey string
	cmd := &cobra.Command{
		Use:   "proc",
		Short: "Show top processes",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := gopsutil.New()
			ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
			defer cancel()

			data, err := c.Processes(ctx, top, sortKey)
			if err != nil {
				return err
			}
			return render.Print(cmd.OutOrStdout(), data)

		},
	}
	cmd.Flags().IntVar(&top, "top", viper.GetInt("process.default_top"), "number of top processes to show")
	cmd.Flags().StringVar(&sortKey, "sort", "cpu", "sort by field: cpu|mem|pid|name")
	return cmd
}
