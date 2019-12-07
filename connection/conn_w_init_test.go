package connection

import (
	"fmt"
	"testing"

	"github.com/jdrivas/gafw/config"
)

func Test_InitialStack(t *testing.T) {
	c := GetCurrentConnection()
	if c.Name != config.DefaultConnectionNameValue {
		t.Errorf("Checking names, got: %s, expected %s", c.Name, config.DefaultConnectionNameValue)
		fmt.Printf("Stack: %#+v\n", currentConnections)
	}
}
func Test_Stack_Growth(t *testing.T) {

	c1 := &Connection{Name: "Test-1"}
	// c2 := &Connection{Name: "Test-2"}
	// c3 := &Connection{Name: "Test-3"}

	expectedLen := 2
	PushCurrentConnection(c1)
	if len(currentConnections) != expectedLen {
		t.Errorf("Pushed one, length of stack is: %d, expect it to be: %d.", len(currentConnections), expectedLen)
		fmt.Printf("Stack: %#+v\n", currentConnections)
	}

}
