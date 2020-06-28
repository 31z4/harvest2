// Package trie implements a radix tree for string prefix analysis.
package trie

import "strings"

// data represents information that a tree edge holds.
type data struct {
	// label is an arbitrary string associated with an edge.
	label string
	// count is the number of times the label was seen when inserting to a tree.
	count uint
}

// edge of a tree.
type edge struct {
	// target points to the node that is connected to this edge.
	target *node
	data
}

// node of a tree.
type node struct {
	// edges connected to this node.
	edges []*edge
}

// Trie implements a radix tree.
// For information about radix trees, see https://en.wikipedia.org/wiki/Radix_tree.
type Trie struct {
	root *node
}

// WalkFn is used when walking a tree.
type WalkFn func(prefix string, count uint)

// New returns empty Trie.
func New() *Trie {
	return &Trie{&node{[]*edge{}}}
}

// shortestLen returns length of the shortest of two strings.
func shortestLen(a, b string) (l int) {
	l = len(a)
	if lb := len(b); lb < l {
		l = lb
	}
	return
}

// longestCommonPrefixLen returns length of the longest common prefix of two strings.
// Note, that this function compares bytes, rather than UTF-8-encoded runes.
// Therefore, it may yield unexpected result for strings encoded in UTF-8.
func longestCommonPrefixLen(a, b string) (l int) {
	s := shortestLen(a, b)
	for l = 0; l < s; l++ {
		if a[l] != b[l] {
			break
		}
	}
	return
}

// appendEdge appends an edge to the node.
func (n *node) appendEdge(e *edge) {
	n.edges = append(n.edges, e)
}

// walk the node recursively calling the specified function for every edge.
func (n node) walk(fn WalkFn, prefixes []string) {
	for _, edge := range n.edges {
		newPrefixes := append(prefixes, edge.label)
		fn(strings.Join(newPrefixes, ""), edge.count)

		if edge.target != nil {
			edge.target.walk(fn, newPrefixes)
		}
	}
}

// split the edge into two edges.
// For example, splitting an edge labeled "test" with index of 2 results in a new edge labeled "st"
// and the original edge labeled "te".
func (e *edge) split(index int) {
	oldLabel := e.label
	oldTarget := e.target

	e.label = oldLabel[:index]
	e.target = new(node)
	e.target.appendEdge(&edge{
		data: data{
			label: oldLabel[index:],
			count: e.count - 1,
		},
		target: oldTarget,
	})
}

// Insert adds a string value to the tree.
func (tree *Trie) Insert(value string) {
	nextEdge := new(edge)
	traverseNode := tree.root
	bytesFound := 0

	for traverseNode != nil {
		nextEdge = nil

		for _, edge := range traverseNode.edges {
			prefixLen := longestCommonPrefixLen(value[bytesFound:], edge.label)
			if prefixLen == 0 {
				continue
			}

			edge.count++
			bytesFound += prefixLen
			nextEdge = edge

			if prefixLen < len(edge.label) {
				edge.split(prefixLen)
			}
			break
		}

		if nextEdge == nil {
			break
		}
		traverseNode = nextEdge.target
	}

	if bytesFound < len(value) {
		if traverseNode == nil {
			traverseNode = new(node)
			nextEdge.target = traverseNode
		}
		traverseNode.appendEdge(&edge{
			data: data{
				label: value[bytesFound:],
				count: 1,
			},
		})
	}
}

// Walk the tree recursively calling the specified function for every edge.
func (tree Trie) Walk(fn WalkFn) {
	tree.root.walk(fn, []string{})
}
