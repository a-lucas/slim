package trie

import (
	"fmt"
	"strings"
)

// Tree defines required functions to convert a tree-like data structure into
// string.
//
// A nil indicates root node.
//
// Since 0.5.1
type Tree interface {

	// Child returns the child of a node and value existence
	Child(node, branch interface{}) interface{}

	// Branches returns the branches for a node.
	Branches(node interface{}) []interface{}

	// NodeID returns the a node id string.
	NodeID(node interface{}) string

	// NodeInfo returns string describing the node.
	NodeInfo(node interface{}) string

	// LeafVal returns the value if the node is leaf.
	LeafVal(node interface{}) (interface{}, bool)
}

// ToString converts a Tree into a multiline string.
//
// A tree node is in form of
//   <income-label>-><node-id>+<step>*<fanout-count>=<value>
// E.g.:
//   000->#000+2*3
//            001->#001+4*2
//                     003->#004+1=0
//                              006->#007+2=1
//                     004->#005+1=2
//                              006->#008+2=3
//            002->#002+3=4
//                     006->#006+2=5
//                              006->#009+2=6
//            003->#003+5=7`[1:]
//
// Since 0.5.1
func ToString(t Tree) string {
	lines := toStrings(t, nil, nil)
	return strings.Join(lines, "\n")
}

// toStrings converts a node and its subtree into multi strings.
//
// Since 0.5.1
func toStrings(t Tree, inbranch, node interface{}) []string {

	line, ind := nodeStr(t, inbranch, node)

	indent := strings.Repeat(" ", ind)

	rst := make([]string, 0, 64)
	rst = append(rst, line)

	for _, b := range t.Branches(node) {
		sub := toStrings(t, b, t.Child(node, b))
		for _, s := range sub {
			rst = append(rst, indent+s)
		}
	}
	return rst
}

// nodeStr returns a string representation of a node, and subtree indent.
//
// Since 0.5.1
func nodeStr(t Tree, inbranch, node interface{}) (string, int) {

	var line string

	if inbranch != nil {
		line += fmt.Sprintf("-%03v->", inbranch)
	}

	nodeid := t.NodeID(node)
	if nodeid != "" {
		line += "#" + nodeid
	}

	indent := len(line)

	line += t.NodeInfo(node)

	brCnt := len(t.Branches(node))
	if brCnt > 1 {
		line += fmt.Sprintf("*%d", brCnt)
	}

	v, isLeaf := t.LeafVal(node)
	if isLeaf {
		line += fmt.Sprintf("=%v", v)
	}
	return line, indent
}
