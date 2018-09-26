// Implementation of an R-Way trie data structure.
//
// A trie has a root node which is the base of the tree.
// Each subsequent node has a letter and children, which are
// nodes that have letter values associated with them.

package zipfs

import "sort"

type node struct {
	val      rune
	term     bool
	depth    int
	meta     interface{}
	mask     uint64
	parent   *node
	children map[rune]*node
}

type trie struct {
	root *node
	size int
}

type byKeys []string

func (a byKeys) Len() int           { return len(a) }
func (a byKeys) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byKeys) Less(i, j int) bool { return len(a[i]) < len(a[j]) }

const nul = 0x0

// Creates a newTrie trie with an initialized root node.
func newTrie() *trie {
	return &trie{
		root: &node{children: make(map[rune]*node), depth: 0},
		size: 0,
	}
}

// Returns the root node for the trie.
func (t *trie) Root() *node {
	return t.root
}

// Adds the key to the trie, including meta data. Meta data
// is stored as `interface{}` and must be type cast by
// the caller.
func (t *trie) Add(key string, meta interface{}) *node {
	t.size++
	runes := []rune(key)
	bitmask := maskruneslice(runes)
	node := t.root
	node.mask |= bitmask
	for i := range runes {
		r := runes[i]
		bitmask = maskruneslice(runes[i:])
		if n, ok := node.children[r]; ok {
			node = n
			node.mask |= bitmask
		} else {
			node = node.NewChild(r, bitmask, nil, false)
		}
	}
	node = node.NewChild(nul, 0, meta, true)
	return node
}

// Finds and returns meta data associated
// with `key`.
func (t *trie) Find(key string) (*node, bool) {
	node := findNode(t.Root(), []rune(key))
	if node == nil {
		return nil, false
	}

	node, ok := node.Children()[nul]
	if !ok || !node.term {
		return nil, false
	}

	return node, true
}

func (t *trie) HasKeysWithPrefix(key string) bool {
	node := findNode(t.Root(), []rune(key))
	return node != nil
}

// Removes a key from the trie, ensuring that
// all bitmasks up to root are appropriately recalculated.
func (t *trie) Remove(key string) {
	var (
		i    int
		rs   = []rune(key)
		node = findNode(t.Root(), []rune(key))
	)

	t.size--
	for n := node.Parent(); n != nil; n = n.Parent() {
		i++
		if len(n.Children()) > 1 {
			r := rs[len(rs)-i]
			n.RemoveChild(r)
			break
		}
	}
}

// Returns all the keys currently stored in the trie.
func (t *trie) Keys() []string {
	return t.PrefixSearch("")
}

// Performs a fuzzy search against the keys in the trie.
func (t trie) FuzzySearch(pre string) []string {
	keys := fuzzycollect(t.Root(), []rune(pre))
	sort.Sort(byKeys(keys))
	return keys
}

// Performs a prefix search against the keys in the trie.
func (t trie) PrefixSearch(pre string) []string {
	node := findNode(t.Root(), []rune(pre))
	if node == nil {
		return nil
	}

	return collect(node)
}

// Creates and returns a pointer to a newTrie child for the node.
func (n *node) NewChild(val rune, bitmask uint64, meta interface{}, term bool) *node {
	node := &node{
		val:      val,
		mask:     bitmask,
		term:     term,
		meta:     meta,
		parent:   n,
		children: make(map[rune]*node),
		depth:    n.depth + 1,
	}
	n.children[val] = node
	n.mask |= bitmask
	return node
}

func (n *node) RemoveChild(r rune) {
	delete(n.children, r)
	for nd := n.parent; nd != nil; nd = nd.parent {
		nd.mask ^= nd.mask
		nd.mask |= uint64(1) << uint64(nd.val-'a')
		for _, c := range nd.children {
			nd.mask |= c.mask
		}
	}
}

// Returns the parent of this node.
func (n node) Parent() *node {
	return n.parent
}

// Returns the meta information of this node.
func (n node) Meta() interface{} {
	return n.meta
}

// Returns the children of this node.
func (n node) Children() map[rune]*node {
	return n.children
}

func (n node) Terminating() bool {
	return n.term
}

func (n node) Val() rune {
	return n.val
}

func (n node) Depth() int {
	return n.depth
}

// Returns a uint64 representing the current
// mask of this node.
func (n node) Mask() uint64 {
	return n.mask
}

func findNode(node *node, runes []rune) *node {
	if node == nil {
		return nil
	}

	if len(runes) == 0 {
		return node
	}

	n, ok := node.Children()[runes[0]]
	if !ok {
		return nil
	}

	var nrunes []rune
	if len(runes) > 1 {
		nrunes = runes[1:]
	} else {
		nrunes = runes[0:0]
	}

	return findNode(n, nrunes)
}

func maskruneslice(rs []rune) uint64 {
	var m uint64
	for _, r := range rs {
		m |= uint64(1) << uint64(r-'a')
	}
	return m
}

func collect(_node *node) []string {
	var (
		keys []string
		n    *node
		i    int
	)
	nodes := []*node{_node}
	for l := len(nodes); l != 0; l = len(nodes) {
		i = l - 1
		n = nodes[i]
		nodes = nodes[:i]
		for _, c := range n.children {
			nodes = append(nodes, c)
		}
		if n.term {
			word := ""
			for p := n.parent; p.depth != 0; p = p.parent {
				word = string(p.val) + word
			}
			keys = append(keys, word)
		}
	}
	return keys
}

type potentialSubtree struct {
	idx  int
	node *node
}

func fuzzycollect(node *node, partial []rune) []string {
	var (
		m    uint64
		i    int
		p    potentialSubtree
		keys []string
	)

	potential := []potentialSubtree{potentialSubtree{node: node, idx: 0}}
	for l := len(potential); l > 0; l = len(potential) {
		i = l - 1
		p = potential[i]
		potential = potential[:i]
		m = maskruneslice(partial[p.idx:])
		if (p.node.mask & m) != m {
			continue
		}

		if p.node.val == partial[p.idx] {
			p.idx++
			if p.idx == len(partial) {
				keys = append(keys, collect(p.node)...)
				continue
			}
		}

		for _, c := range p.node.children {
			potential = append(potential, potentialSubtree{node: c, idx: p.idx})
		}
	}
	return keys
}