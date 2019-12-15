package cmd

import (
	"fmt"
	"net/http"

	"github.com/jdrivas/conman"
	t "github.com/jdrivas/termtext"
	"github.com/spf13/viper"
)

// Shim for t.httpDisplay, with timing display.
func httpDisplay(se *conman.SideEffect, resp *http.Response, err error) {
	if !viper.GetBool(t.JSONDisplayKey) {
		if se.ElapsedTime.Milliseconds() < 1000 {
			fmt.Printf(t.Title("Command took %d milliseconds\n", se.ElapsedTime.Milliseconds()))
		} else {
			fmt.Printf(t.Title("Command took %4g seconds\n", se.ElapsedTime.Seconds()))

		}
	}
	t.HTTPDisplay(resp, err)
}
