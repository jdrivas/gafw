// Copyright © 2018 David Rivas

package cmd

import (
	"fmt"
	"os"

	"github.com/jdrivas/gafw/config"
	t "github.com/jdrivas/gafw/term"
	"github.com/spf13/cobra"

	// "github.com/spf13/pflag"
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
			DoInteractive()
		},
	}
	// Add the commands to the rootCmd node (e.g. http get /users).
	if mode != interactive {
		rootCmd.AddCommand(interactiveCmd)
	}

	// httpCmd = &cobra.Command{
	// 	Use:   "http",
	// 	Short: "Use HTTP verbs.",
	// 	Long:  "Send requests to the service with HTTP verbs and arguments.",
	// }
	// rootCmd.AddCommand(httpCmd)

	// buildHTTP(mode)
}

//
// Flag and config file init.
//

var (
	// tokenFV, hubURLFV                                  string
	// authClientIDFV, authClientSecretFV, authRedirectFV string

	verbose, debug bool
)

// InitCmd is designed to be used from Main - ordering is important here so can't just execute whenever.
// Should only be called once.
func InitCmd() {
	fmt.Printf("%s\n", t.Title("InitCmd"))

	// Root is created here, rather than in build root to deal with command line
	// flags in interactive mode.
	// Any root command flags set on the original command line should persist
	// to _each_ interactive command. They can  be explicitly overridden if needed.
	rootCmd = &cobra.Command{
		Use:   fmt.Sprintf("%s <command> [<args>]", config.AppName),
		Short: "Talk to a forest server.",
		Long:  "A tool for working with a forest server.",
	}

	// We initFlags here because rootCmd.Execute() will parse the command line
	// and error out if flags are present that are not defined.
	initFlags()
	config.InitConfig()
	fmt.Printf("%s\n", t.SubTitle("Debug is: %t", viper.GetBool(config.DebugKey)))

	// Each time we run a command we should do cobraInit()
	// cobra.OnInitialize(cobraInit)
	fmt.Printf("%s\n", t.Title("InitCmd - exit"))

}

func initFlags() {
	fmt.Printf("%s\n", t.Title("initFlags"))

	// Rest flags to start
	rootCmd.ResetFlags()

	// Flags available to everyone.
	rootCmd.PersistentFlags().StringVar(&config.ConfigFileName, config.ConfigFlagKey, "", fmt.Sprintf("config file location. (default is %s{yaml,json,toml}", config.ConfigFileRoot))

	rootCmd.PersistentFlags().BoolVarP(&verbose, config.VerboseFlagKey, "v", false, "Describe what is happening as its happening.")
	viper.BindPFlag(config.VerboseFlagKey, rootCmd.PersistentFlags().Lookup(config.VerboseFlagKey))

	rootCmd.PersistentFlags().BoolVarP(&debug, config.DebugFlagKey, "d", false, "Describe details about what's happening.")
	viper.BindPFlag(config.DebugFlagKey, rootCmd.PersistentFlags().Lookup(config.DebugFlagKey))

	fmt.Printf("%s\n", t.SubTitle("Debug is: %t", viper.GetBool(config.DebugKey)))

	// fmt.Printf("%s\n", t.Title("End of InitFlags - Viper dump"))
	// viper.Debug()
	// fmt.Printf("%s\n", t.Title("End of InitFlags - End of Viper dump"))

}

// This should be called AFTER the config file has been read.
func initConnectionWithFlags() {
	fmt.Printf("%s\n", t.Title("InitConnectionWithFlags"))

	// Do the normal config file default
	// initConnections()

}

// Intended to be executed once before each commend.
// This happens after the commands line has been parsed
// but before any CMDs have been executed.
func cobraInit() {
	fmt.Printf("%s\n", t.Title("cobraInit"))
	fmt.Printf("%s\n", t.SubTitle("Debug is: %t", viper.GetBool(config.DebugKey)))

	// config.InitConfig()
	initFlags()
	initConnectionWithFlags()
	config.InitConfig()
	fmt.Printf("%s\n", t.Title("cobraInit - exit"))

}
