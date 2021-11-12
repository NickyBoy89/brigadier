package brigadier

import (
	"fmt"
	"io"

	"github.com/Tnze/go-mc/net/packet"
)

// NodeType is the type of command node that a node is
type NodeType byte

func (nt NodeType) String() string {
	switch nt {
	case ROOT:
		return "ROOT"
	case LITERAL:
		return "LITERAL"
	case ARGUMENT:
		return "ARGUMENT"
	default:
		return "UNUSED"
	}
}

const (
	ROOT NodeType = iota
	LITERAL
	ARGUMENT
	UNUSED
)

// Node is a brigadier structure that defines a notchian command and how it connects to other commands.
// A node can be one of three types, ROOT, LITERAL, or ARGUMENT, and can be serialized as a go-mc packet.
// A Node's zero value is a ROOT node, with no options set.
// A reference for a node's structure can be found here: https://wiki.vg/Command_Data
type Node struct {
	Flags           int8          // Flags for the node, contains type, and other configuration
	Children        []int32       // Other command nodes that make up the rest of the commmand
	RedirectNode    int32         // Optional, only if 0x08 of flags is set
	Name            string        // Optional, only for ARGUMENT and LITERAL nodes
	Parser          string        // Optional, only for ARGUMENT nodes
	Properties      []interface{} // Currently unimplemented, only for ARGUMENT nodes
	SuggestionsType string        // Optional, only if 0x10 of flags is set
}

// NewNode creates a new brigadier node with the specified options set
func NewNode(nodeType NodeType, isExecutable, hasRedirect, hasSuggestions bool) Node {
	node := Node{}

	// Set the bit flags
	node.Flags |= int8(nodeType)
	if isExecutable {
		node.Flags |= 0x04
	}
	if hasRedirect {
		node.Flags |= 0x08
	}
	if hasSuggestions {
		node.Flags |= 0x10
	}
	return node
}

func (n Node) String() string {
	return fmt.Sprintf("CommandNode type: %v, flags: %.5x", NodeType(n.Flags&0x03), n.Flags)
}

// WriteTo defines how a node will be encoded as a packet
func (n Node) WriteTo(w io.Writer) (int64, error) {
	var total int64
	count, err := packet.Byte(n.Flags).WriteTo(w)
	if err != nil {
		panic(err)
	}
	total += count
	count, err = packet.VarInt(len(n.Children)).WriteTo(w)
	if err != nil {
		panic(err)
	}
	total += count
	convertedChildren := make([]packet.VarInt, len(n.Children))
	for i, child := range n.Children {
		convertedChildren[i] = packet.VarInt(child)
	}
	count, err = packet.Ary{Ary: convertedChildren, Len: len(n.Children)}.WriteTo(w)
	if err != nil {
		panic(err)
	}
	total += count
	// Optional int
	count, err = packet.Opt{Has: n.Flags&0x08 == 1, Field: packet.VarInt(n.RedirectNode)}.WriteTo(w)
	if err != nil {
		panic(err)
	}
	total += count
	// Optional string
	count, err = packet.Opt{Has: n.Flags&0x02 == 1 || n.Flags&0x01 == 1, Field: packet.String(n.Name)}.WriteTo(w)
	if err != nil {
		panic(err)
	}
	total += count
	// Optional string
	count, err = packet.Opt{Has: n.Flags&0x02 == 1, Field: packet.String(n.Parser)}.WriteTo(w)
	if err != nil {
		panic(err)
	}
	total += count
	/*
		count, err = packet.ByteArray(n.Properties).WriteTo(w)
		if err != nil {
			panic(err)
		}
		total += count
	*/
	count, err = packet.Opt{Has: n.Flags&0x10 == 1, Field: packet.String(n.SuggestionsType)}.WriteTo(w)
	if err != nil {
		panic(err)
	}
	total += count
	return total, nil
}

// ReadFrom defines how a node will be decoded from a packet
func (n *Node) ReadFrom(r io.Reader) (int64, error) {
	var (
		Flags         packet.Byte
		ChildrenCount packet.VarInt
		Children      = []packet.VarInt{}
		RedirectNode  packet.VarInt
		Name          packet.String
		Parser        packet.Identifier
		//Properties      []interface{}
		SuggestionsType packet.Identifier
	)

	var total int64
	count, err := Flags.ReadFrom(r)
	if err != nil {
		panic(err)
	}
	total += count
	count, err = ChildrenCount.ReadFrom(r)
	if err != nil {
		panic(err)
	}
	total += count
	children := packet.Ary{Ary: &Children, Len: ChildrenCount}
	count, err = children.ReadFrom(r)
	if err != nil {
		panic(err)
	}
	total += count
	// Optional int
	count, err = packet.Opt{Has: Flags&0x08 == 1, Field: RedirectNode}.ReadFrom(r)
	if err != nil {
		panic(err)
	}
	total += count
	// Optional string
	count, err = packet.Opt{Has: Flags&0x01 == 1 || Flags&0x02 == 1, Field: Name}.ReadFrom(r)
	if err != nil {
		panic(err)
	}
	total += count
	// Optional string
	count, err = packet.Opt{Has: Flags&0x11 == 1, Field: Parser}.ReadFrom(r)
	if err != nil {
		panic(err)
	}
	total += count
	/*
		count, err = packet.ByteArray(n.Properties).ReadFrom(r)
		if err != nil {
			panic(err)
		}
		total += count
	*/
	count, err = packet.Opt{Has: Flags&0x11 == 1, Field: SuggestionsType}.ReadFrom(r)
	if err != nil {
		panic(err)
	}
	total += count
	return total, nil
}
