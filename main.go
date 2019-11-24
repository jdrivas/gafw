// Copyright © 2018 David Rivas

package main

import (
	"github.com/jdrivas/gcli/cmd"
	"github.com/jdrivas/gcli/config"
	"github.com/jdrivas/gcli/term"
)

func main() {

	config.InitConfig()
	cmd.InitCmd()
	// fmt.Printf("Configuration:\n")
	// for k, v := range viper.AllSettings() {
	// 	fmt.Printf("%s: %v\n", k, v)
	// }

	term.InitTerm()

	// Off you go ....
	cmd.Execute()
}
