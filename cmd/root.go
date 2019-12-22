// Copyright © 2018 David Rivas

package cmd

import (
	"fmt"
	"os"

	connection "github.com/jdrivas/conman"
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

	// Second level state commands

	showCmd.AddCommand(&cobra.Command{
		Use:   "flags",
		Short: "view flags",
		Long:  "Display the flags for this appliation and their current settings.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", t.SubTitle("Flags are:"))
			var flags []*pflag.FlagSet
			// flags = append(flags, rootCmd.PersistentFlags()) // it appears persistent appear in local.
			flags = append(flags, rootCmd.LocalFlags())
			w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 2, ' ', 0)
			fmt.Fprintf(w, "%s", t.Title("Name\tShort\tValue\tDefValue\tChanged\n"))
			for _, fs := range flags {
				fs.VisitAll(func(f *pflag.Flag) {
					fmt.Fprintf(w, "%s\n", t.SubTitle("%s\t%s\t%s\t%s\t%t",
						f.Name, f.Shorthand, f.Value.String(), f.DefValue, f.Changed))
				})
			}
			w.Flush()
			if viper.GetBool(config.DebugFlagKey) {
				fmt.Printf("viper Debug key is set.\n")
			}

		},
	})

	showCmd.AddCommand(&cobra.Command{
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
	fmt.Printf(t.Title("cmd/root/init()\n")) // We can't bracket this with config.Debug as viper won't be set yet.

	config.AppName = "gafw"
	cobra.OnInitialize(doCobraOnInit)

	rootCmd = &cobra.Command{
		Use:   fmt.Sprintf("%s <command> [<args>]", config.AppName),
		Short: "Talk to a forest server.",
		Long:  "A tool for working with a forest server.",
	}

	// Wants to happen ahead of cobra initialization.
	// Flags are parsed before the cobra.OnInitialization call.
	// This is why we call resetEnvironment which calls initFlags(), in
	// resetEnvironment before we call at the end of
	// the interactive process loop doICommand()
	initFlags()

	fmt.Printf(t.Title("cmd/root/init() - exit\n\n"))

}

var firstCobraInit = true
var setFlags map[string]*pflag.Flag

// Called before the command has been executed
// but after the
func doCobraOnInit() {
	// yes, yes we could do real tracing ....
	if config.Debug() {
		fmt.Printf(t.Title("doCobraOnInit()\n"))
	}

	if firstCobraInit {
		if config.Debug() {
			fmt.Printf(t.Alert("First Init\n"))
		}

		// Get all the set flags.
		flags := rootCmd.PersistentFlags()
		setFlags = make(map[string]*pflag.Flag)
		flags.VisitAll(func(f *pflag.Flag) {
			fmt.Printf("Flag: %s, Changed: %t\n", f.Name, f.Changed)
			if f.Changed {
				setFlags[f.Name] = f
			}
		})

		if config.Debug() {
			fmt.Printf("Found %d set flags\n", len(setFlags))
		}
		firstCobraInit = false
	} else {
		if config.Debug() {
			fmt.Printf(t.Success("Not first init\n"))
		}

	}

	// Always do config first, as many inits depend on it.
	config.InitConfig()
	t.InitTerm()
	connection.InitConnections()

	if config.Debug() {
		fmt.Printf(t.Title("doCobraOnInit() - exit\n\n"))
	}
}

var (
	debug, verbose, json bool
	screenProfileName    string
)

const (
	jsonFlagKey          = "json"
	screenProfileFlagKey = "screen"
)

// Flag behaviour in interactive mode.
// Command line flags set on program invocation (e.g. gafw --vernbose interactive) will stick throughout the
// interactive session. (Of course in this particular case debug can be reset at the prompt).
//
// Flags used at the interactive prompt will only have effect on the command executed, and not remain durable
// for future commands.
// Thus at the interactive prompt:
//   gafw [con-name https://foo.bar.com]: http -c test-conn get /
// will do a get on whateer is defined by test-conn and then go back to using https://foo.bar.com

// These flags are meant to be provide on program invocation and remain durable throughout
// an interactive session. That is, they don't get reset prior to each interactive prompt.
// e.g.
//      gafw --screen termDarkDefault interactive
// this will set the screen to termDeftaulInteractive for the duration of the interactive session.
//
// whereas:
//     gafw -c https://foo.bar.com interactive

//
// Something needs to initializze all the flags prior to config.InitConfig() which initializes viper
// and reads in any configuration file.
//
// resetENvironemnt() defined in interactive.go needs to call something that only
// resets flags that were not set on the program invocation command line.
//
// doCobraInit is the only thing that is guranteed to get us in before commands are executed
// and after the command line has been read. Thus, in doCobraInit we must:
// * Check to see if this is the first time through.
// * If so store flag state for flags that were explicitly set on the application command line.
// * if not resetFlags, and reply the ones that we stored the first time through.
//
// In addition we can consider removing the Cobra falgs/command reset stuff out of the interactive code
// and putting it in doCobraOnInit():
//  * rootCmd.ResetFlags()
//  * rootCmd.ResetCommands()
//  * buildRoot(interactive)
//  * initFlags()
// and leave resetEnvirontment to just adding the commands to the tree for interactive mode: e.g. rootCmd.AddCommand(exitCmd)

// This is pulled out of the general init because
// we want to refer to it in interactive mode, where we first rootCmd.ResetFlags() to start
// from scratch and then call this to reset the flags.
// Note, tried the idea of just running through the existing values and resting to
// defualt (or other values as the case may be). This "works" but has the side effect
// of updating the Changed state in the flag.
// This allows a one time set of the flags from the interactive command line.
// Now, once you create a flag, you can't redefine it so, you have to eraase them first.
// this is why resetEnvironment does that.
func initFlags() {
	fmt.Printf(t.Title("initFlags\n"))

	// Flags available to everyone.
	rootCmd.PersistentFlags().StringVar(&config.ConfigFileName, config.ConfigFlagKey, "",
		fmt.Sprintf("config file location. (default is %s{yaml,json,toml}", config.ConfigFileRoot))

	defaultVerbose := false
	if f, ok := setFlags[config.VerboseFlagKey]; ok {
		if f.Value.String() == "true" {
			defaultVerbose = true
		}
	}
	rootCmd.PersistentFlags().BoolVarP(&verbose, config.VerboseFlagKey, "v", defaultVerbose,
		"Describe what is happening as its happening.")
	viper.BindPFlag(config.VerboseFlagKey, rootCmd.PersistentFlags().Lookup(config.VerboseFlagKey))

	defaultDebug := false
	if f, ok := setFlags[config.DebugFlagKey]; ok {
		if f.Value.String() == "true" {
			defaultDebug = true
		}
	}
	rootCmd.PersistentFlags().BoolVarP(&debug, config.DebugFlagKey, "d", defaultDebug,
		"Describe details about what's happening.")
	viper.BindPFlag(config.DebugFlagKey, rootCmd.PersistentFlags().Lookup(config.DebugFlagKey))

	defaultConnection := ""
	if f, ok := setFlags[connection.ConnectionFlagKey]; ok {
		defaultConnection = f.Value.String()
	}
	rootCmd.PersistentFlags().StringVarP(&connection.ConnectionFlagValue, connection.ConnectionFlagKey, "c", defaultConnection,
		"Use the named connection (names defined in config file)")

	defaultJSON := false
	if f, ok := setFlags[jsonFlagKey]; ok {
		if f.Value.String() == "true" {
			defaultJSON = true
		}
	}
	rootCmd.PersistentFlags().BoolVarP(&json, jsonFlagKey, "j", defaultJSON,
		"Print output in unencumbred JSON for easy scripting.")
	viper.BindPFlag(t.JSONDisplayKey, rootCmd.PersistentFlags().Lookup(jsonFlagKey))

	// ScreenProfile
	defaultScreenProfile := t.ScreenNoColorDefaultKey
	if f, ok := setFlags[screenProfileFlagKey]; ok {
		defaultScreenProfile = f.Value.String()
	}
	rootCmd.PersistentFlags().StringVarP(&screenProfileName, screenProfileFlagKey, "s", defaultScreenProfile,
		"Set the screen profile for output (e.g. colors etc).")
	viper.BindPFlag(t.ScreenProfileKey, rootCmd.PersistentFlags().Lookup(screenProfileFlagKey))

	fmt.Printf(t.Title("initFlags -- done\n"))
}
