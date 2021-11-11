package brigadier

import (
	"reflect"
	"testing"

	"github.com/Tnze/go-mc/net/packet"
)

func TestSetNodeFields(t *testing.T) {
	testNode := NewNode(ARGUMENT, true, true, true)
	if testNode.Flags != 0b11110 {
		t.Errorf("Node had byte flag of %.5b, expected %.5b", testNode.Flags, 0b11110)
	}
}

func TestBasicDecodeEncodeNode(t *testing.T) {
	testNode := NewNode(ROOT, false, false, false)
	testPacket := packet.Marshal(
		0x00,
		testNode,
	)

	var (
		OutputNode Node
	)

	if err := testPacket.Scan(&OutputNode); err != nil {
		t.Errorf("Error decoding packet: %v", err)
	}

	t.Logf("Sent: %v", testNode)
	t.Logf("Received: %v", OutputNode)

	if !reflect.DeepEqual(testNode, OutputNode) {
		t.Errorf("Nodes were not the same: original %v, and decoded %v", testNode, OutputNode)
	}
}
