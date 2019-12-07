package connection

import (
	"fmt"
	"testing"
)

func Test_Stack_PushPop(t *testing.T) {

	c := &Connection{
		Name: "test-1",
	}
	PushCurrentConnection(c)

	if len(currentConnections) != 1 {
		t.Errorf("Stack is wrong size. Should be 1, got %d", len(currentConnections))
		fmt.Printf("Stack: %v\n", currentConnections)
	}

	nc := PopCurrentConnection()
	if len(currentConnections) != 0 {
		t.Errorf("Stack is wrong size. Should be 0, got %d", len(currentConnections))
		fmt.Printf("Stack: %v\n", currentConnections)
	}
	if c != nc {
		t.Errorf("Popped the wrong elmeent. \nt\tShould be %v\n\tgot %v", c, nc)
		fmt.Printf("Stack: %v\n", currentConnections)
	}

}

func Test_Stack_Bounds(t *testing.T) {

	nameBase := "Test-"
	var testSize = 10
	for i := 0; i < testSize; i++ {
		c := &Connection{
			Name: fmt.Sprintf("%s%d", nameBase, i),
		}
		PushCurrentConnection(c)
	}

	if len(currentConnections) != testSize {
		t.Errorf("Pushed %d connections, but length was %d.", testSize, len(currentConnections))
		fmt.Printf("Stack: %v\n", currentConnections)
	}

	for i := testSize - 1; i >= 0; i-- {
		expectName := fmt.Sprintf("%s%d", nameBase, i)
		c := PopCurrentConnection()
		if c.Name != expectName {
			t.Errorf("Checking names, got: %s, expected %s", c.Name, expectName)
			fmt.Printf("Stack: %#+v\n", currentConnections)
		}
	}

}
