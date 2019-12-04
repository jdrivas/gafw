package connection

import (
	"fmt"

	"github.com/jdrivas/gafw/config"
	"github.com/jdrivas/gafw/term"
	"github.com/spf13/viper"
)

// Connection contains information for connecting to a service endpoint.
type Connection struct {
	Name       string
	ServiceURL string
	AuthToken  string
	Headers    map[string]string
}

// ConnectionList is a colleciton of connections.
type ConnectionList []*Connection

var currentConnection *Connection

// GetCurrentConnection is the primary interface for getting the connection to use.
func GetCurrentConnection() *Connection {
	return currentConnection
}

// SetCurrentConnection sets the connection in use.
func SetCurrentConnection(conn *Connection) {
	currentConnection = conn
}

// GetAllConnections returns a list of known connections
func GetAllConnections() ConnectionList {
	return getAllConnectionsFromConfig()
}

// GetConnection by name
func GetConnection(name string) (*Connection, bool) {
	return getConnectionFromConfig(name)
}

// SetConnection sets the current connection by name
func SetConnection(name string) (err error) {
	conn, ok := GetConnection(name)
	if ok {
		SetCurrentConnection(conn)
	} else {
		err = fmt.Errorf("couldn't find connection \"%s\"", name)
	}
	return err
}

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

// Get a connection from the config file:
// connections:
//    name:
//        serviceURL: https://foo.bar.com
//        authToken: XXXX-YYYYY-XXXXX
//        Headers:
//            X-User-Id: 1234557-1237AC
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
func InitConnections() {
	if config.Debug() {
		fmt.Printf("Initializing Connections\n")
	}

	// Current conenction should be durable during interactive mode
	// reset it to the default ...
	var conn *Connection
	if currentConnection == nil {
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
				if config.Debug() {
					fmt.Printf("Using a 'broken' default connection.\n")
				}
				// ... As a last resort set up a broken empty connection.
				// We won't panic here as we can set it during interactive
				// mode and it will otherwise error.
				conn = &Connection{
					Name:       config.DefaultConnectionNameValue,
					ServiceURL: defaultServiceURL,
				}
			}
		}
		if config.Debug() {
			fmt.Printf("Using connection: %s[%s]\n", conn.Name, conn.ServiceURL)
		}
		// lastConnection = conn
		SetCurrentConnection(conn)
	}
	// 	// or if we've just changed it for one command, reset it to previous.
	// } else if getCurrentConnection().Name == updatedConnectionName {
	// 	conn = lastConnection
	// 	SetCurrentConnection(conn)
	// }
}
