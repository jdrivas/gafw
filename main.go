// Copyright © 2018 David Rivas

package main

import (
	"github.com/jdrivas/gafw/cmd"
)

func main() {

	// config.InitConfig()
	// cmd.InitCmd() // Set up the initial commands.
	// fmt.Printf("Configuration:\n")
	// for k, v := range viper.AllSettings() {
	// 	fmt.Printf("%s: %v\n", k, v)
	// }

	// Off you go ....
	cmd.Execute()
}
