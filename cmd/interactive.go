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
	"github.com/spf13/viper"
)

var (

	// Type exit instead of just control-d, Note: We actually os.exit() here.
	// Which eans no post-processing of any kind including simply falling through
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
			viper.Set(config.VerboseKey, !viper.GetBool(config.VerboseKey))
			vs := "Off"
			if viper.GetBool(config.VerboseKey) {
				vs = "On"
			}
			fmt.Printf("Verbose is %s\n", vs)
		},
	}

	debugCmd = &cobra.Command{
		Use:     "debug",
		Aliases: []string{"d"},
		Short:   "Toggle debug mode and print status.",
		Long:    "Toggle debug, verbose will print out detailed status as its happening.",
		Run: func(cmd *cobra.Command, args []string) {
			viper.Set(config.DebugKey, !viper.GetBool(config.DebugKey))
			vs := "Off"
			if viper.GetBool(config.DebugKey) {
				vs = "On"
			}
			fmt.Printf("Debug is %s\n", vs)
		},
	}
)

// Each time through the loop we rebuild the command tree
//  and reinitialize the flags.
func resetEnvironment() {
	if config.Debug() {
		fmt.Printf("%s\n", t.Title("interactive.resetEnvironment()"))
	}

	// We must rebuild the flags and the command
	// each time through.
	// Remember
	rootCmd.ResetFlags()
	rootCmd.ResetCommands()
	buildRoot(interactive)
	initFlags()
	rootCmd.AddCommand(exitCmd)
	rootCmd.AddCommand(verboseCmd)
	rootCmd.AddCommand(debugCmd)
	connection.ResetConnection()

	// initialize the flags on the tree
	// initFlags()
	// config.InitConfig()
	// This is only here to reset the prmopt
	// TODO: The connection handling logicis is
	// a disaster. Fix it.
	// initConnectionWithFlags()
	if config.Debug() {
		fmt.Printf("%s\n", t.Title("interactive.resetEnvironment() - exit"))
	}

}

// Parse the line and execute the command
func doICommand(line string) (err error) {

	rootCmd.SetArgs(strings.Fields(line)) // Don't use strings.Split - it won't eat white space.
	err = rootCmd.Execute()

	resetEnvironment()
	return err
}

func promptLoop(process func(string) error) (err error) {

	// Set up for the first itme through.
	resetEnvironment()

	for moreCommands := true; moreCommands; {
		conn := connection.GetCurrentConnection()
		serviceURL := conn.ServiceURL
		connName := conn.Name
		// token := conn.getSafeToken(true, false)
		token := ""
		spacer := ""
		if token != "" {
			spacer = " "
		}
		status := statusDisplay()
		prompt := fmt.Sprintf("%s [%s%s %s]: ", t.Title(config.AppName), t.Info(status), t.Highlight(connName), t.SubTitle("%s%s%s", serviceURL, spacer, token))
		// prompt := fmt.Sprintf("%s [%s%s]: ", t.Title(config.AppName), t.Info(status), t.SubTitle("context"))

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

// DoInteractive sets up a readline loop that reads and executes comands.
func DoInteractive() {
	readline.SetHistoryPath(fmt.Sprintf("./%s", config.HistoryFile))
	xICommand := func(line string) (err error) { return doICommand(line) }
	err := promptLoop(xICommand)
	if err != nil {
		fmt.Printf("Error exiting prompter: %s\n", t.Fail(err.Error()))
	}
}
