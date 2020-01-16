package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	connection "github.com/jdrivas/conman"
	t "github.com/jdrivas/termtext"
	config "github.com/jdrivas/vconfig"
	"github.com/spf13/cobra"
)

var (

	// Type exit instead of just control-d, Note: We actually os.exit() here.
	// Which means no post-processing of any kind if there was any by simply falling through
	// to the orignial Execute command.
	// if this is a problem, move the definition of the promptLoop moreCommands up
	// to module scope and set it to false in the Run function directly below.
	exitCmd = &cobra.Command{
		Use:     "exit",
		Aliases: []string{"quit"},
		Short:   "Exit from the application",
		Long:    "Stop reading input lines and terminate the application.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\nGoodbye and thank you.\n")
			os.Exit(0)
		},
	}

	verboseCmd = &cobra.Command{
		Use:     "verbose",
		Aliases: []string{"v"},
		Short:   "Toggle verbose mode and print status.",
		Long:    "Toggle verbose, verbose will print out detailed status as its happening.",
		Run: func(cmd *cobra.Command, args []string) {
			config.ToggleVerbose()
			// boundFlags.remove(config.VerboseFlagKey)
			vs := "Off"
			if config.Verbose() {
				vs = "On"
			}
			fmt.Printf("Verbose is %s\n", vs)
		},
	}

	debugCmd = &cobra.Command{
		Use:     "debug",
		Aliases: []string{"d"},
		Short:   "Toggle debug mode and print status.",
		Long:    "Toggle debug, debug will print out detailed status as its happening.",
		Run: func(cmd *cobra.Command, args []string) {
			// fmt.Printf("set flags: %p, \n", &boundFlags)
			config.ToggleDebug()
			// boundFlags.remove(config.DebugFlagKey) // do this to remove the effect of application flags once a set has been done.
			// fmt.Printf("boundFlags: %p\n", &boundFlags)
			vs := "Off"
			if config.Verbose() {
				vs = "On"
			}
			fmt.Printf("Debug is %s\n", vs)
		},
	}
)

// Add the above into the tree off of rootCmd
func addInteractiveCommands() {
	rootCmd.AddCommand(exitCmd)
	rootCmd.AddCommand(verboseCmd)
	rootCmd.AddCommand(debugCmd)
}

// DoInteractive sets up a readline loop that reads and executes comands.
// This is the entrypoint for application command line interactive command.
func DoInteractive() {

	addInteractiveCommands()

	readline.SetHistoryPath(fmt.Sprintf("./%s", config.HistoryFile))

	xICommand := func(line string) (err error) { return doICommand(line) }
	err := promptLoop(xICommand)
	if err != nil {
		fmt.Printf("Error exiting prompter: %s\n", t.Fail(err.Error()))
	}
}

// Feed the line to Cobra at the root command.
// Then execute rootCmd.
func doICommand(line string) (err error) {

	// rootCmd.ResetCommands()
	// buildRoot(interactive)
	// addInteractiveCommands()

	args := strings.Fields(line) // Don't use strings.Split - it won't eat white space.
	rootCmd.ParseFlags(args)
	rootCmd.SetArgs(args)
	err = rootCmd.Execute()

	return err
}

// Build prompt, readline, manage history, until it's time to stop.
func promptLoop(process func(string) error) (err error) {

	for moreCommands := true; moreCommands; {
		serviceURL := ""
		connName := ""
		if conn, err := connection.GetCurrentConnection(); err == nil {
			serviceURL = conn.ServiceURL
			connName = conn.Name
		}
		// token := conn.getSafeToken(true, false)
		token := ""
		spacer := ""
		if token != "" {
			spacer = " "
		}
		status := statusDisplay()
		prompt := fmt.Sprintf("%s [%s%s %s]: ",
			t.Title(config.AppName), t.Info(status), t.Highlight(connName),
			t.SubTitle("%s%s%s", serviceURL, spacer, token))

		if config.Debug() {
			fmt.Println() // add a stanza mark between the spew.
		}
		line, err := readline.Line(prompt)
		if err == io.EOF {
			moreCommands = false
		} else if err != nil {
			fmt.Printf("Readline Error: %s\n", t.Fail(err.Error()))
		} else {
			readline.AddHistory(line)
			err = process(line)
			if err == io.EOF {
				moreCommands = false
			}
		}
	}
	return nil
}

// Yes, I'm sure there's some kind of []rune
// thing to do here instead.
func statusDisplay() (s string) {
	if config.Verbose() {
		s = fmt.Sprintf("%s%s", s, "v")
	}
	if config.Debug() {
		s = fmt.Sprintf("%s%s", s, "d")
	}
	if len(s) > 0 {
		s = fmt.Sprintf("%s%s", s, " ")
	}
	return s
}
