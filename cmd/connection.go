package cmd

import (
	"github.com/jdrivas/gafw/connection"
	"github.com/jdrivas/gafw/term"
	"github.com/spf13/cobra"
)

var connectionAliases = []string{"conn", "con"}

func buildConnection(runMode) {

	listCmd.AddCommand(&cobra.Command{
		Use:     "connections [flags]",
		Aliases: connectionAliases,
		Short:   "list connections available for service",
		Long:    "Display a description of the connections available to sue to send HTTP commands.",
		Args:    cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			conns := connection.GetAllConnections()
			term.List(conns, nil, nil)
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
				term.List(connection.GetAllConnections(), nil, err)
			}
		},
	})

}
