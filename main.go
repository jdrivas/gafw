// Copyright © 2018 David Rivas

package main

import (
	"github.com/jdrivas/gafw/cmd"
	"github.com/jdrivas/gafw/term"
)

func main() {

	// config.InitConfig()
	// cmd.InitCmd() // Set up the initial commands.
	// fmt.Printf("Configuration:\n")
	// for k, v := range viper.AllSettings() {
	// 	fmt.Printf("%s: %v\n", k, v)
	// }

	term.InitTerm()

	// Off you go ....
	cmd.Execute()
}
