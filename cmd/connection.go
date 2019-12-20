package cmd

import (
	connection "github.com/jdrivas/conman"
	t "github.com/jdrivas/termtext"
	"github.com/spf13/cobra"
)

var connectionAliases = []string{"conns", "conn", "con"}

func buildConnection(runMode) {

	listCmd.AddCommand(&cobra.Command{
		Use:     "connections [flags]",
		Aliases: connectionAliases,
		Short:   "list connections available for service",
		Long:    "Display a short description of the connections available to use to send HTTP commands.",
		Args:    cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			conns := connection.GetAllConnections()
			t.List(conns, nil, nil)
		},
	})

	describeCmd.AddCommand(&cobra.Command{
		Use:     "connection [flags] <connection-name> ...",
		Aliases: connectionAliases,
		Short:   "Details about a service connection",
		Long:    "Display details about a connection or connections that are available to send HTTP commands.",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			conns := connection.GetAllConnections()
			var fconns connection.ConnectionList
			for _, cn := range args {
				if c := conns.FindConnection(cn); c != nil {
					fconns = append(fconns, c)
				}
			}
			t.Describe(fconns, nil, nil)
		},
	})

	setCmd.AddCommand(&cobra.Command{
		Use:     "connection <connection-name>",
		Aliases: connectionAliases,
		Short:   "Use the named connection.",
		Long:    "Sets the service connection to the named connection.",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := connection.SetConnection(args[0])
			if err != nil {
				t.List(connection.GetAllConnections(), nil, err)
			}
		},
	})

}
