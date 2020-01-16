// Copyright © 2018 David Rivas

package cmd

import (
	"fmt"
	"os"

	"github.com/jdrivas/gafw/version"
	t "github.com/jdrivas/termtext"
	config "github.com/jdrivas/vconfig"
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
	listCmd, describeCmd    *cobra.Command
	setCmd, showCmd         *cobra.Command
)

// This is pulled out specially, because for interactive
// Root of the command hierarcy. All commands reference one of these.
func buildRoot(mode runMode) {

	// First level commands

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

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List objects",
		Long:  "Short description of a collection of objects.",
	}
	rootCmd.AddCommand(listCmd)

	describeCmd = &cobra.Command{
		Use:   "describe",
		Short: "Details about an object",
		Long:  "Detailed description of a an object.",
	}
	rootCmd.AddCommand(describeCmd)

	setCmd = &cobra.Command{
		Use:   "set",
		Short: "Set set",
		Long:  "Set a value or values to the application state or on an object.",
	}
	rootCmd.AddCommand(setCmd)

	showCmd = &cobra.Command{
		Use:   "show",
		Short: "Describe state",
		Long:  "Short description of applcation state state.",
	}
	rootCmd.AddCommand(showCmd)

	httpCmd = &cobra.Command{
		Use:   "http",
		Short: "Use HTTP verbs.",
		Long:  "Send requests to the service with HTTP verbs and arguments.",
	}
	rootCmd.AddCommand(httpCmd)

	//
	// Second level state commands
	//

	// Config/Flags

	showCmd.AddCommand(&cobra.Command{
		Use:   "flags",
		Short: "view flags",
		Long:  "Display the flags for this appliation and their current settings.",
		Run: func(cmd *cobra.Command, args []string) {
			flags := rootCmd.LocalFlags()
			printFlagSet(flags)
		},
	})

	showCmd.AddCommand(&cobra.Command{
		Use:   "app-flags",
		Short: "view application-flags",
		Long:  "Display the flags as set on the application innvocation from the command line.",
		Run: func(cmd *cobra.Command, args []string) {
			for i, bf := range config.GetBindFlags() {
				fmt.Printf("%d: %#v\n", i, bf)
			}
		},
	})

	showCmd.AddCommand(&cobra.Command{
		Use:     "config",
		Aliases: []string{"cnofig"},
		Short:   "view configuration",
		Long:    "Display the configuration information for this application as set by file, evnironment, flags",
		Run: func(cmd *cobra.Command, args []string) {
			printConfig()
		},
	})

	// Terminal Profile

	setCmd.AddCommand(&cobra.Command{
		Use: "screen",
		// Aliases: []string{""},
		Short: "set terminal profile",
		Long:  "Set the color scheme and other atttibutes of the terminal profile",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// boundFlags.remove(screenProfileFlagKey)
			config.Set(t.ScreenProfileKey, args[0])
			// viper.Set(t.ScreenProfileKey, args[0])
		},
	})

	showCmd.AddCommand(&cobra.Command{
		Use: "screen",
		// Aliases: []string{""},
		Short: "view terminal profile",
		Long:  "View name and attributes of the terminal profile that is current set.",
		Run: func(cmd *cobra.Command, args []string) {
			p := viper.GetString(t.ScreenProfileKey)
			fmt.Printf("Screen profile is \"%s\": %s %s %s %s ",
				p, t.Title("Title"), t.SubTitle("SubTitle"), t.Text("Text"), t.Highlight("Highlight"))
			fmt.Printf("%s %s %s %s\n",
				t.Success("Success"), t.Warn("Warn"), t.Fail("Fail"), t.Alert("Alert"))

		},
	})

	// Build out sub menus.
	buildHTTP(mode)
	buildConnection(mode)
}

func displayFlags(fs *pflag.FlagSet) {
	w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 2, ' ', 0)
	fmt.Fprintf(w, "%s", t.Title("Name\tShort\tValue\tDefValue\tChanged\n"))
	fs.VisitAll(func(f *pflag.Flag) {
		fmt.Fprintf(w, "%s\n", t.SubTitle("%s\t%s\t%s\t%s\t%t",
			f.Name, f.Shorthand, f.Value.String(), f.DefValue, f.Changed))
	})
	w.Flush()

}

//
// Flag and config file init.
//

func init() {
	viper.AutomaticEnv()
	if config.Debug() {
		t.Pef()
		defer t.Pxf()
	}

	config.AppName = "gafw"
	cobra.OnInitialize(doCobraOnInit)

	rootCmd = &cobra.Command{
		Use:               fmt.Sprintf("%s <command> [<args>]", config.AppName),
		Short:             "Talk to a forest server.",
		Long:              "A tool for working with a forest server.",
		PersistentPreRun:  rootPre,
		PersistentPostRun: rootPost,
	}

	// Wants to happen ahead before flags are actually parsed of course.
	// Flags are parsed before the cobra.OnInitialization call.
	initFlags()

}
