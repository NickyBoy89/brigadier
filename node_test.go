package brigadier

import "testing"

func TestSetNodeFields(t *testing.T) {
	testNode := NewNode(ARGUMENT, true, true, true)
	if testNode.Flags != 0b11110 {
		t.Errorf("Node had byte flag of %.5b, expected %.5b", testNode.Flags, 0b11110)
	}
}
