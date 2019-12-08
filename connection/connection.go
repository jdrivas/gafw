package connection

import (
	"fmt"
	"os"
	"sort"

	"github.com/jdrivas/gafw/config"
	"github.com/jdrivas/gafw/term"
	t "github.com/jdrivas/gafw/term"
	"github.com/juju/ansiterm"
	"github.com/spf13/viper"
)

//
// Public API
//

// Connection contains information for connecting to a service endpoint.
type Connection struct {
	Name       string
	ServiceURL string
	AuthToken  string
	Headers    map[string]string
}

// ConnectionList is a colleciton of connections.
type ConnectionList []*Connection

// Sort by name
type byName ConnectionList

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[j], a[i] = a[i], a[j] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// Manage a stack of connections
var currentConnections = make([]*Connection, 0)

// GetCurrentConnection is the primary interface for getting the connection to use.
func GetCurrentConnection() *Connection {
	if len(currentConnections) == 0 {
		return nil
	}

	// Get the top of the stack
	return currentConnections[len(currentConnections)-1]
}

// GetConnection by name (from configuration).
func GetConnection(name string) (*Connection, bool) {
	return getConnectionFromConfig(name)
}

// SetConnection sets the current connection by name
// That is, replace the stop of the stack with the onen given.
func SetConnection(name string) (err error) {
	conn, ok := GetConnection(name)
	if ok {
		PopCurrentConnection()
		PushCurrentConnection(conn)
	} else {
		err = fmt.Errorf("couldn't find connection \"%s\"", name)
	}
	return err
}

// GetAllConnections returns a list of known connections
func GetAllConnections() ConnectionList {
	return getAllConnectionsFromConfig()
}

//
// The current connection is implemented with a stack
// to facilitate changing the connection via a flag.

// PushCurrentConnection  adds the connection to the stack, setting the new current connection.
func PushCurrentConnection(conn *Connection) {
	// fmt.Printf("Push - Size of conn stack: %d\n", len(currentConnections))
	// Don't push a dupilicate on the stack.
	if len(currentConnections) == 0 || GetCurrentConnection().Name != conn.Name {
		currentConnections = append(currentConnections, conn)
	}
	// fmt.Printf("Push done - Size of conn stack: %d\n", len(currentConnections))
}

// PopCurrentConnection returns the top of the stack and sets the new
// current connection to the next topmost stack item.
// Returns nil on empty stack.
func PopCurrentConnection() (c *Connection) {
	// fmt.Printf("Pop - Size of conn stack: %d\n", len(currentConnections))
	if len(currentConnections) > 0 {
		currentConnections, c = currentConnections[:len(currentConnections)-1], currentConnections[len(currentConnections)-1]
	}
	// fmt.Printf("Pop Done - Size of conn stack: %d\n", len(currentConnections))
	return c
}

// List displpays the list of connections and notes the current one.
func (conns ConnectionList) List() {
	sort.Sort(byName(conns))
	if len(conns) > 0 {
		currentName := GetCurrentConnection().Name
		w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
		fmt.Fprintf(w, "%s\n", t.Title("\tName\tURL"))
		for _, c := range conns {
			name := term.Text(c.Name)
			current := ""
			if c.Name == currentName {
				name = term.Highlight("%s", c.Name)
				current = term.Highlight("%s", "*")
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", current, name, t.Text("%s", c.ServiceURL))
		}
		w.Flush()
	} else {
		fmt.Printf("%s\n", t.Title("There were no connections."))
	}
}

// Private API
// Read in the config to get all the named connections
func getAllConnectionsFromConfig() (conns ConnectionList) {
	connectionsMap := viper.GetStringMap(config.ConnectionsKey) // map[string]interface{}
	for name := range connectionsMap {
		conn, ok := getConnectionFromConfig(name)
		if ok {
			conns = append(conns, conn)
		} else {
			fmt.Printf(term.Error(fmt.Errorf("couldn't create a config for connection \"%s\"", name)))
		}
	}
	return conns
}

func getConnectionFromConfig(name string) (conn *Connection, ok bool) {
	connKey := fmt.Sprintf("%s.%s", config.ConnectionsKey, name)
	if viper.IsSet(connKey) {
		conn = &Connection{
			Name:       name,
			ServiceURL: viper.GetString(fmt.Sprintf("%s.%s", connKey, config.ServiceURLKey)),
			AuthToken:  viper.GetString(fmt.Sprintf("%s.%s", connKey, config.AuthTokenKey)),
			Headers:    viper.GetStringMapString(fmt.Sprintf("%s.%s", connKey, config.HeadersKey)),
		}
		ok = true
	}
	return conn, ok
}

// initConnections sets up the first current Connection,
// initializes the ShowTokens state, and should be called whenever the Viper config file gets reloaded.
// Since we need at least a URL to break and/or let us know that no token has been set.
const defaultServiceURL = "http://127.0.0.1:80"

// InitConnections initializes a default connection. Needs to happen after we've read in the viper configuration file.
// TODO: It's probably best if init is idempotent.
func InitConnections() {
	if config.Debug() {
		fmt.Printf("Initializing Connections\n")
	}
	fmt.Printf("Initializing Connections\n")

	// Current conenction should be durable during interactive mode
	// reset it to the default ...
	var conn *Connection
	if len(currentConnections) == 0 {
		if config.Debug() {
			fmt.Printf("No current Connection.\n")
		}
		// If there is a connection named default, use it ....
		var ok bool
		conn, ok = GetConnection(config.DefaultConnectionNameValue)
		if !ok {
			// .. Otherwise, see if there is a _name_ of a defined connection to use as default ...
			defaultName := viper.GetString(config.DefaultConnectionNameKey)
			conn, ok = GetConnection(defaultName)
			if !ok {

				// ... next look for _any_ defined connections.
				// Rather than pick a random connection (maps don't have a determined order.
				// and we get connections from the config file as a map), pick the first lexographic one.
				conns := getAllConnectionsFromConfig()
				if len(conns) > 0 {
					sort.Sort(byName(conns))
					conn = conns[0]
				} else {
					// ... As a last resort set up a broken empty connection.
					// We won't panic here as we can set it during interactive
					// mode and it will otherwise error.
					if config.Debug() {
						fmt.Printf("Using a 'broken' default connection.\n")
					}
					conn = &Connection{
						Name:       config.DefaultConnectionNameValue,
						ServiceURL: defaultServiceURL,
					}
				}
			}
		}
		if config.Debug() {
			fmt.Printf("Using connection: %s[%s]\n", conn.Name, conn.ServiceURL)
		}

		PushCurrentConnection(conn)
	}
}
