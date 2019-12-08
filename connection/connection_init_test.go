package connection

import (
	"fmt"
	"testing"

	"github.com/jdrivas/gafw/config"
	"github.com/spf13/viper"
)

// Test to make sure that when InitConnections is called
// the correct connection is set.
func Test_InitialStack(t *testing.T) {
	InitConnections()
	setConnectionConfig("test-1", "htpps://localhost")
	c := GetCurrentConnection()
	if c.Name != config.DefaultConnectionNameValue {
		t.Errorf("Checking names, got: %s, expected %s", c.Name, config.DefaultConnectionNameValue)
		fmt.Printf("Stack: %#+v\n", currentConnections)
	}

	resetConnection()
}

// the rest of the tests will mock out viper and
// assum an empty stack.
func resetConnection() {
	viper.Reset()
	currentConnections = make([]*Connection, 0)
}
