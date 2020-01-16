package cmd

import (
	"fmt"
	"strings"

	connection "github.com/jdrivas/conman"
	config "github.com/jdrivas/vconfig"
	"github.com/spf13/cobra"
)

func buildHTTP(mode runMode) {
	// HTTP Util
	// TODO: Consider validating the HTTP verbs.
	httpCmd.AddCommand(&cobra.Command{
		Use:                   "send [flags] <method> <command> [<json-string> ....]",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"SEND"},
		Short:                 "HTTP <method> <command> to the service.",
		Long: `Sends an HTTP <method> <command> to the current service endpoint.
<method> is an HTTP verb (e.g. "GET")

All of the args following <command> are caputred as a single json 
string and placed in the body of the request, 
with the ContentType header set to application/json.`,
		Example: fmt.Sprintf(" %s http send post /groups/test/users {\"name\": \"admin\", \"users\": [\"david\"]}", config.AppName),
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if conn, err := connection.GetCurrentConnection(); err == nil {
				if len(args) == 2 {
					httpDisplay(conn.Send(strings.ToUpper(args[0]), args[1], nil, nil))
				} else {
					httpDisplay(conn.Send(strings.ToUpper(args[0]), args[1], strings.Join(args[2:], " "), nil))
				}
			}
		},
	})

	httpCmd.AddCommand(&cobra.Command{
		Use:                   "get [flags]  <command>",
		Aliases:               []string{"GET"},
		DisableFlagsInUseLine: true,
		Short:                 "HTTP GET <command> to service.",
		Args:                  cobra.MinimumNArgs(1),
		Long:                  " Sends an HTTP GET <command> to the service endpoint.",
		Example:               fmt.Sprintf("%s http get /users", config.AppName),
		Run: func(cmd *cobra.Command, args []string) {
			if conn, err := connection.GetCurrentConnection(); err == nil {
				httpDisplay(conn.Get(args[0], nil))
			}
		},
	})

	httpCmd.AddCommand(&cobra.Command{
		Use:                   "post [flags] <command> [<json-string> ....]",
		Aliases:               []string{"POST"},
		DisableFlagsInUseLine: true,
		Short:                 "HTTP POST <command> <body> to service.",
		Long: `Sends an HTTP POST <command> <body> to the service endpoint.  

All of the args follwing <command> are caputred as a single json 
string and placed in the body of the request, 
with the ContentType header set to application/json.`,
		Example: fmt.Sprintf("%s http post /groups/test/users {\"name\": \"admin\", \"users\": [\"david\"]}", config.AppName),
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if conn, err := connection.GetCurrentConnection(); err == nil {
				if len(args) > 1 {
					httpDisplay(conn.Post(args[0], strings.Join(args[1:], " "), nil))
				} else {
					httpDisplay(conn.Post(args[0], nil, nil))
				}
			}
		},
	})

	httpCmd.AddCommand(&cobra.Command{
		Use:     "delete [flags] <command> [<json-string> ....]",
		Aliases: []string{"DELETE"},
		Short:   "HTTP DELETE <command> <body> to service.",
		Long: `Sends an HTTP DELETE <command> <body>to the service endpoint.  

All of the args following <command> are caputred as a single json 
string and placed in the body of the request, 
with the ContentType header set to application/json.`,
		Example: fmt.Sprintf("%s http delete /groups/test/users {\"name\": \"admin\", \"users\": [\"david\"]}", config.AppName),
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if conn, err := connection.GetCurrentConnection(); err == nil {
				httpDisplay(conn.Delete(args[0], nil, nil))
			}
		},
	})

}
