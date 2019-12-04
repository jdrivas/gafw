// Copyright © 2018 David Rivas

package cmd

import (
	"fmt"
	"os"

	"github.com/jdrivas/gafw/config"
	"github.com/jdrivas/gafw/connection"
	"github.com/jdrivas/gafw/term"
	t "github.com/jdrivas/gafw/term"
	"github.com/jdrivas/gafw/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/juju/ansiterm"
	"github.com/spf13/viper"
)

// it gets run before each line is parsed.
// runMode allows us to add or remove commands
// as necessary for interadtive use
type runMode int

const (
	interactive runMode = iota + 1
	commandline
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	buildRoot(commandline)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// commands
var (
	rootCmd, interactiveCmd *cobra.Command
	httpCmd                 *cobra.Command
)

// This is pulled out specially, because for interactive
// Root of the command hierarcy. All commands reference one of these.
func buildRoot(mode runMode) {

	interactiveCmd = &cobra.Command{
		Use:   "interactive",
		Short: "Interactive mode",
		Long:  "Runs a command line interpreter with sematnics to make session use easy.",
		Run: func(cmd *cobra.Command, args []string) {
			if config.Debug() {
				fmt.Printf("Processing interactive.\n")
			}
			DoInteractive()
			if config.Debug() {
				fmt.Printf("interactive done.\n")
			}
		},
	}
	// Add the commands to the rootCmd node (e.g. http get /users).
	if mode != interactive {
		rootCmd.AddCommand(interactiveCmd)
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version.",
		Long:  "Every program needs a version, this shows you what the value is.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", version.Version)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "flags",
		Short: "view flags",
		Long:  "Display the flags for this appliation and their current settings.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", t.SubTitle("Flags are:"))
			fs := rootCmd.PersistentFlags()
			w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 2, ' ', 0)
			fmt.Fprintf(w, "%s", t.Title("Name\tShort\tValue\tDefValue\tChanged\n"))
			fs.VisitAll(func(f *pflag.Flag) {
				fmt.Fprintf(w, "%s\n", t.SubTitle("%s\t%s\t%s\t%s\t%t",
					f.Name, f.Shorthand, f.Value.String(), f.DefValue, f.Changed))
			})
			w.Flush()
			if viper.GetBool(config.DebugFlagKey) {
				fmt.Printf("viper Debug key is set.\n")
			}

		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "config",
		Short: "view configuration",
		Long:  "Display the configuration information for this application as set by file, evnironment, flags",
		Run: func(cmd *cobra.Command, args []string) {
			settings := viper.AllSettings()
			fmt.Printf("%s\n", t.SubTitle("Settings are:"))
			w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 2, ' ', 0)
			fmt.Fprintf(w, "%s", t.Title("Name\tValue\n"))
			for k := range settings {
				fmt.Fprintf(w, "%s\t%s\n", t.Title("%s", k), t.SubTitle("%#+v", settings[k]))
			}
			w.Flush()
			if viper.GetBool(config.DebugFlagKey) {
				fmt.Printf("viper Debug key is set.\n")
			}

		},
	})

	httpCmd = &cobra.Command{
		Use:   "http",
		Short: "Use HTTP verbs.",
		Long:  "Send requests to the service with HTTP verbs and arguments.",
	}
	rootCmd.AddCommand(httpCmd)
	buildHTTP(mode)
}

//
// Flag and config file init.
//

// InitCmd is designed to be used from Main - ordering is important here so can't just execute whenever.
// Should only be called once.
func init() {
	// fmt.Printf("%s\n", t.Title("InitCmd"))  // We can't bracket this with config.Debug as viper won't be set yet.

	cobra.OnInitialize(doCobraOnInit)

	rootCmd = &cobra.Command{
		Use:   fmt.Sprintf("%s <command> [<args>]", config.AppName),
		Short: "Talk to a forest server.",
		Long:  "A tool for working with a forest server.",
	}

	// Wants to happen ahead of cobra initialization.
	// Flags are parsed before the cobra.OnInitialization call.
	initFlags()

	// fmt.Printf("%s\n", t.Title("InitCmd - exit"))

}

func doCobraOnInit() {
	// yes, yes we could do real tracing ....
	if config.Debug() {
		fmt.Printf("%s\n", t.Title("doCobraOnInit()"))
	}
	config.InitConfig()
	term.InitTerm()
	connection.InitConnections()

	if config.Debug() {
		fmt.Printf("%s\n", t.Title("doCobraOnInit() - exit"))
	}
}

var (
	debug, verbose bool
)

// This is pulled out of the general init because
// we want to refer to it in interactive mode, where we first rootCmd.ResetFlags() to start
// from scratch and then call this to reset the flags.
// This allows a one time set of the flags from the interactive command line.
func initFlags() {
	// fmt.Printf("%s\n", t.Title("initFlags"))

	// Flags available to everyone.
	rootCmd.PersistentFlags().StringVar(&config.ConfigFileName, config.ConfigFlagKey, "", fmt.Sprintf("config file location. (default is %s{yaml,json,toml}", config.ConfigFileRoot))

	rootCmd.PersistentFlags().BoolVarP(&verbose, config.VerboseFlagKey, "v", false, "Describe what is happening as its happening.")
	viper.BindPFlag(config.VerboseFlagKey, rootCmd.PersistentFlags().Lookup(config.VerboseFlagKey))

	rootCmd.PersistentFlags().BoolVarP(&debug, config.DebugFlagKey, "d", false, "Describe details about what's happening.")
	viper.BindPFlag(config.DebugFlagKey, rootCmd.PersistentFlags().Lookup(config.DebugFlagKey))

	// fmt.Printf("%s\n", t.Title("initFlags -- done"))
}
