package cmd

import (
	"fmt"

	connection "github.com/jdrivas/conman"
	t "github.com/jdrivas/termtext"
	config "github.com/jdrivas/vconfig"
	"github.com/spf13/cobra"
)

/*
Configuration and Interactive Mode

Command line flags set at program invocation (e.g. gafw --screen termDarkDefault interactive) will
stick throughout the interactive session. That is, they don't get reset prior to each interactive prompt.
e.g.
     jdr> gafw --screen termDarkDefault interactive
will set the screen-profile to termDarkDeftault for the duration of the interactive session.

Flags used at the interactive prompt will only have effect on the command executed, and not remain durable
for future commands.
Thus at the interactive prompt:
  gafw [con-name https://foo.bar.com]: http --screen termNoColorDefault get /
will display, for this command only, without color.

There are some flags which will set a mode and there are commands to reset the mode, eg. verbose.
The behavior there is that if the flag is set at the command line it will behave durablly until
a "verbose" command is issued, which toggles the mode and so would rset verbose mode durablly (e.g verbose
would be off). The interactive prompt use of the flag remains applicable only for the life of the single
interactive command.
Thus:
jdr> gafw --verbose interactive                                       # Sets verbose on
gafw [v con-name https://foo.bar.com]: http  get /                    # verbose still on
gafw [v con-name https://foo.bar.com]: verbose                        # verbose off
gafw [on-name https://foo.bar.com]: http  get /                       # verbose still off
gafw [on-name https://foo.bar.com]: http  --verbose get /             # verbose on only for this command
gafw [on-name https://foo.bar.com]: http  get /						  # verbose remains off.


How Initialization Works with Cobra and Viper.

The usual order of appoliation operation with Cobra and Viiper is:
1. Build a cobra command tree.
2. Define flags and their bindings to confugration keys for lookup in Viper.
3. Configure and initialize Viper: Point at the config file, bind to the environment variables,
   read in the configuration file, set up watch on the condig file  if you want it.

Since flags are attached to commands in Cobra, 1 must preceede 2. It seems that 2 preceeds
3 here primarily to provide a flag for locating the configuration file.

Important note: Viper doesn't bind flags, environment variables, and configuration file values
until lookup. Configuration is immediate. So, for example, if an environment variable is set,
this can be seen once the environment has been bound to Viper. Thus, we bind Viper to the
environment right away in cmd/root.go:init(), to get Debug through the environment even before
we've read in a configuration file or parsed a flag.

Viper Precedence

Viper implements predence of operation in binding a value to a viper key in the following order:

1. Explicit call to Set
2. Flag
3. Environment
4. Config file
5. Key/Value store
6. Default


How to Itegrate With an Interactive Command Loop.

If you are using a variable that is going to be managed by Viper. Keep the state in Viper, dont create a
go variable that you then attempt to mirror. Mirror variables just cause pain and issues of
synchronizaing state.

Thus vconfig supports: vconfig.Debug() and manages the swtate inside of viper. Thus
vconfig.Debug(), vconfig.ToggleDebug() all use viper.Get/Set on DebugKey.

Issues
1. Everytime we go through the interactive loop we want to apply configuration as if we're starting
   fresh, at least with respect to flags (well, almost). Deleting all of the previous flags and then
   add the flags back on the commnand tree is easiest (and perhaps it's the only way to get what
   we want). Note, you can't, for instance, simply get the flags and reset values. In particular,
   this doesn't work with the flag.Changed value.

   The tear down and build up is implemented in cmd.reset() below. To support the basic behavior
   reset needs to happen before the command line is parsed. The only entry point before the command line
   is parsed is after a command has been exexuted (and so prior to the next command). This could be
   implemented in the interactive loop itself, but cobra provides a couple of command entry points.
   The PersistentPostRun command is used in this case.

2. Package update and initialization. To support initialization of other tools that want to also leverage
   Viper configuration it would be handy to send an event saying paramaters have updated.
   Relevant events are:
          * Each time we reset, paramaters are essentially sent back to a default state.
          * Whenever comnand line processing is done and interactive command line arguments are parsed
            the flags may change the configuration values.
          * Whenever an evnironment variable changes.
          * Whenver the config file is re-read in.

   While we can get hooks to all of these except env var change, it easier just to note that there are
   really only a two places where change will show:
          * Right before a command is executed
          * Right after we've reset but before the command prompt is printed. The Command prompt often
            displays state that may rely on configuraiton.

   So, we update in both pre and post functions. The pre case gets the command line flags, the post case
   resets to the default state for the command line prompt. This last point is particularly important
   for resetting to default case after a flag has altered state for just the duration of a single command.
   The command line prompt is outside the scope of the single command duration.

3. Integration and priority of applciation command line.
   With the above fixes most things work as expected. A problem is that the behavior of application
   command line flags do not remain durable across interactive command line invocations. Essentially,
   this use of flags introduces another level of precedence into the system, in between environment variable and
   interactive command line flag.

   Actually there is another level of precedence that is controlled by any set commands on
   a state (eg. verbose or connections). This looks like:

   The new desired precedence then,from low to high:

               6. application default              # Lowest precedence
               5. config file
               4. environment variable
               3. application command line flag
               2. variable Set/Toggle commands
			   1. interactive command line flag    # Highest


	One issue with the above is that this conflicts with the Viper precedence of Set over Flags.

	Thoughts on this last issue:
		The only thing I can think of is that we manage the binding of flags to viper variables
		directly rather than with viper.

		Thus instead of doing a viper.BindPFlag or otherwise, we handle that our self in config.
			vconfig.BindPFlag() ....
		I believe then that the implication is that in the cobra machinary we will have to install
		one or more of cobra.OnInitialize(), and rootCmd.PersistentPreRun, rootCmd.PeristentPostrun
		to deal with flag management and actually doing the flag bind to variables. Note, I think this
		still means that we support the other automatic binding from vipe (env, config file etc), there
		would be little use in still using viper if not.

	In Sum:

	We need to do the following things:
	* Capture all flags that are bound to viper variables.
	* Handle flag binding whenever a flag set is parsed, distinquishing between the application
	  command line flag updates and interactive command line flag updates.
	* To support application defaults, note when a Set has occured to prioritize that
	  over the application default.

   One way of doing this is:
   1. At flag initailization, capture viper/flag binding.
   2. At application startup captiure and maintain any values set by the application commandline flags.
   3. Any time one of those variables is Set, with a command, update the value in the captured state.
   4. Apply the state prior to command execution and after comamand line parsing, and immediately after
      execution and prior to command prompt display.


So:

1. Applcation command line flags. Capture the application command line flags. This state can be grabed in
   OnCobraInit. This only happens the first time OnCobraInit is called, so a first-time sentienel is
   placed in the OnCobraInit function.

   This is implemented as:
   typedef BoundFlag struct {
		Flag *pflag.Flag
		BindKey string
		Value interface{}          // Value set explicityly or by app flags.
   }
   typedef boundFlags map[string]BoundFlag
   Indexed on the Flag.Name (the key used for the flag)

2. Flag Set up and Reset. Flags get set up as a part of initializaiton. This is done in initFlags().
   The initialization needs to happen before Cobra parses the command line, so the first time this
   is called is in cmd/root.go:init(). Flags are both created (defaults set etc) and then bound to a
   viper variable here with vconfig.bindFlag(bindKey string, flag *pflag.Flag)

   In addition each time through the interactive loop, flags should be torn down with rootCmd.ResetFlags(),
   and then initFlags() should be called to buid them back up again(including bindings) This is to support reparsing for each
   command line (it's possible, we can get rid of this with rootCmd.ParseFlags(), rootCmd.SetArgs()).

   The full reset is captured in a function called reset().
   Then this is called in the PostRun function on rootCmd

       rootCmd.ResetFlags()
       initFlags()

3. On the frist run through OnCobraInit, all the Changed flags (those set on the command line)
   values are captured: the Value on the appropriate BoundFlag is set to the value of the flag.
   For those that had not been changed, capture the current value in viper to Value. This durability
   of the value from the initial command line.

4. Set values not managed by Cobra/Viper integration. This is really the business end of implementing the desired
   behavior for the application  command line flags.

   a. Each time Command is about to Execute (either corba.OnInitialize() or rootCmd.PPreRun)
		 i.  If an interactive command line flag has been set for this, set the bound value to the flag value.
		ii.  Else, set the bound value to the SetValue.
   b. In reset(), after initFlags, restore the Set value to each Viper variable. This undoes any interactive
      command line updates.

5. 	Setting bound variables. To set a bound variable, you need to use vconfig.Set(key string, value interface{}.
	This updates the binding to keep the set value	and update enforce the desired presidence and then calls
	viper.Set( ... )

6. Module initialization. This happens in moduleInit, which is called in the PreRun function off rootCmd.
   Examples calls include t.InitTerm(), connection.InitConnections(). This serves as an envent mechanism
   to let modules know that they may want to re-evaluate viper variables due to state change.

*/

// Package Update
//

// packageFuncs is the list of functions we want to call reread viper variables when set/reset.
var packageFuncs = []func(){t.InitTerm, connection.InitConnections}

func moduleInit() {
	if config.Debug() {
		t.Pef()
		defer t.Pxf()
	}
	for _, f := range packageFuncs {
		f()
	}
}

// cobra.OnInitialize registers a function that is called "everytime a command's Execute method is called".
// Praticaly this is like a PersistentPreRun on root but without having to bother with the commnand tree.
var firstCobraInit = true

// Called before the command has been executed
// but after the
func doCobraOnInit() {
	// yes, yes we could do real tracing ....
	if config.Debug() {
		t.Pef()
		defer t.Pxf()
	}

	if firstCobraInit {
		if config.Debug() {
			fmt.Printf(t.Alert("First Init\n"))
		}

		// Get all the set flags from the application (not interactive) command line ...
		config.UpdateChangedFlags()
		// and apply the results of all bound flags to viper
		config.Apply()

		// Set up general configuration (read config file.)
		// We do this here because we want to be able to use
		// the command line flag to modify the configuration file.
		config.InitConfig()
		firstCobraInit = false
	} else {
		if config.Debug() {
			fmt.Printf(t.Alert("After First Init\n"))
		}
		// Pick up flags from the interactive line, and
		// apply them, without updating the bind variables.
		// rootPost will apply BindVariables on another pass.
		config.ApplyFromFlags(rootCmd.PersistentFlags())
	}

}

// Called before a command executes but after flags have been parsed.
// Useful to send a message to packages that variables have been updated.
func rootPre(cmd *cobra.Command, args []string) {
	if config.Debug() {
		t.Pef()
		defer t.Pxf()
	}
	moduleInit()
}

// Called after a command has executed, but before the command line prompt
// has printed.
// Useful to reset values to those before the comnand line flags had been parsed for this
// command.
func rootPost(cmd *cobra.Command, args []string) {
	if config.Debug() {
		t.Pef()
		defer t.Pxf()
	}

	reset()        // Reset the environment for another pass through
	config.Apply() // Reapply apply those set on the command line.
	moduleInit()   // inform interested parties
}

// Reset the flags and bindings.
func reset() {
	if config.Debug() {
		t.Pef()
		defer t.Pxf()
	}

	rootCmd.ResetFlags() // Literally erases the flags from the tree.
	initFlags()
}

// Initialize Flags
//

var (
	debugFlag, verboseFlag, jsonFlag bool
	screenProfileFlag                string
	connectionFlag                   string
)

const (
	configFlagKey        = "configfile"
	debugFlagKey         = "debug"
	verboseFlagKey       = "verbose"
	jsonFlagKey          = "json"
	screenProfileFlagKey = "screen"
	connectionFlagKey    = "connection"
)

// Create flags and bind them to  viper variables.
func initFlags() {
	if config.Debug() {
		t.Pef()
		defer t.Pxf()
	}

	// config file
	rootCmd.PersistentFlags().StringVar(&config.ConfigFileName, configFlagKey, "",
		fmt.Sprintf("config file location. (default is %s{yaml,json,toml}", config.ConfigFileRoot))

	// Verbose
	defaultVerbose := false
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, verboseFlagKey, "v",
		defaultVerbose, "Describe what is happening as its happening.")
	config.Bind(config.VerboseKey, rootCmd.PersistentFlags().Lookup(verboseFlagKey))

	// Debug
	defaultDebug := false
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, debugFlagKey, "d",
		defaultDebug, "Describe details about what's happening.")
	config.Bind(config.DebugKey, rootCmd.PersistentFlags().Lookup(debugFlagKey))

	// Connection
	defaultConnection := ""
	rootCmd.PersistentFlags().StringVarP(&connectionFlag, connectionFlagKey, "c",
		defaultConnection, "Use the named connection (names defined in config file)")
	config.Bind(connection.DefaultConnectionNameKey, rootCmd.PersistentFlags().Lookup(connectionFlagKey))

	// JSON
	defaultJSON := false
	rootCmd.PersistentFlags().BoolVarP(&jsonFlag, jsonFlagKey, "j",
		defaultJSON, "Print output in unencumbred JSON for easy scripting.")
	config.Bind(t.JSONDisplayKey, rootCmd.PersistentFlags().Lookup(jsonFlagKey))

	// ScreenProfile
	defaultScreenProfile := t.ScreenNoColorDefaultKey
	rootCmd.PersistentFlags().StringVarP(&screenProfileFlag, screenProfileFlagKey, "s",
		defaultScreenProfile, "Set the screen profile for output (e.g. colors etc).")
	config.Bind(t.ScreenProfileKey, rootCmd.PersistentFlags().Lookup(screenProfileFlagKey))
}
