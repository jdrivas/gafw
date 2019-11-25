package config

import (
	"fmt"
	"os"

	t "github.com/jdrivas/gafw/term"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	AppName = "gafw"
)

var (
	ConfigFileName string
	ConfigFileRoot = fmt.Sprintf("%s", AppName)
	HistoryFile    = fmt.Sprintf(".%s_history", AppName)
)

/*
* Keys to look up values in Viper configuration.
 */

// TODO: Formalize the yaml file structure, or at least document it here.

// YAML Variables which show up in viper, but managed here.
const (
	DebugKey         = "debug"            // bool
	VerboseKey       = "verbose"          // bool
	ScreenProfileKey = t.ScreenProfileKey // this avoids the circular reference to
)

// Flags These are the long form flag values for command line flags.
const (
	ConfigFlagKey       = "config"
	hubURLFlagKey       = "hub-url"
	TokenFlagKey        = "token"
	AuthRedirectFlagKey = "auth-redirect-url"
	ClientIDFlagKey     = "client-id"
	ClientSecretFlagKey = "client-secret"
	VerboseFlagKey      = "verbose"
	DebugFlagKey        = "debug"
)

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {

	fmt.Printf("%s\n", t.Title("InitConfig"))

	// Fin a config file
	if ConfigFileName != "" {
		viper.SetConfigFile(ConfigFileName)
	} else {
		viper.SetConfigName(ConfigFileRoot)

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra_test" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// Read in the config file.
	if err := viper.ReadInConfig(); err == nil {
		if Debug() {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	} else {
		fmt.Printf("Error loading config file: %s - %v\n", viper.ConfigFileUsed(), err)
	}
	fmt.Printf("%s\n", t.SubTitle("Debug is: %t", viper.GetBool(DebugKey)))
	fmt.Printf("%s\n", t.Title("InitConfig - exit"))

}
