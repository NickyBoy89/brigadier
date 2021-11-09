package brigadier

// NodeType is the type of command node that a node is
type NodeType byte

const (
	ROOT NodeType = iota
	LITERAL
	ARGUMENT
	UNUSED
)

// Node is a brigadier structure that defines a notchian command
// and how it connects to other commands
// Reference: https://wiki.vg/Command_Data
type Node struct {
	Flags           byte
	ChildrenCount   int
	Children        []int
	RedirectNode    *int
	Name            *string
	Parser          *string
	Properties      []interface{}
	SuggestionsType *string
}

// NewNode creates a new brigadier node with the specified options set
func NewNode(nodeType NodeType, isExecutable, hasRedirect, hasSuggestions bool) *Node {
	node := new(Node)

	// Set the bit flags
	node.Flags |= byte(nodeType)
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
