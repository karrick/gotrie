package gotrie

import (
	"strings"
)

// PrefixTrie is a prefix tree, also known as a digital tree, as described by
// https://en.wikipedia.org/wiki/Trie with a one byte radix.
//
//     package main
//
//     import (
//         "bufio"
//         "fmt"
//         "os"
//
//         "github.com/karrick/gotrie"
//     )
//
//     func main() {
//         // build a new Trie from standard input lines
//         t := gotrie.NewPrefixTrie()
//         scanner := bufio.NewScanner(os.Stdin)
//         for scanner.Scan() {
//             t.Insert(scanner.Text(), struct{}{})
//         }
//
//         if err := scanner.Err(); err != nil {
//             fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
//             os.Exit(1)
//         }
//
//         // Enumerate through Trie in sorted order
//         for t.Scan() {
//             fmt.Println(t.Text())
//         }
//     }
type PrefixTrie struct {
	// root represents the root node of the tree, and is associated with the
	// empty string.
	root *pnode

	// bookmarks are used while enumerating trie contents during scanning.
	bookmarks []*pbookmark
}

// NewPrefixTrie returns a new prefix trie.
func NewPrefixTrie() *PrefixTrie {
	return &PrefixTrie{root: new(pnode)}
}

// pbookmark is used while enumerating a prefix trie contents during scanning.
type pbookmark struct {
	n *pnode
	k uint16
}

// pnode is a node in a prefix trie.
type pnode struct {
	children [256]*pnode
	value    interface{}
	valid    bool
}

// Find locates the specified key and returns its respective value, along with a
// boolean which is true when the key was found.
func (t *PrefixTrie) Find(key string) (interface{}, bool) {
	n := t.root
	for _, k := range []byte(key) {
		c := n.children[k]
		if c == nil {
			return nil, false
		}
		n = c
	}
	return n.value, n.valid
}

// Insert stores the key-value pair in the Trie, overwriting an existing value
// if key was stored before.
func (t *PrefixTrie) Insert(key string, value interface{}) {
	n := t.root
	keyb := []byte(key)

	indexLastByte := -1

	for i, k := range keyb {
		c := n.children[k]
		if c == nil {
			indexLastByte = i
			break
		}
		n = c
	}

	// append new nodes for the remaining bytes, if any
	if indexLastByte != -1 {
		for _, k := range keyb[indexLastByte:] {
			c := new(pnode)
			n.children[k] = c
			n = c
		}
	}

	n.value = value
	n.valid = true
}

// Scan locates the next key-value pair in the Trie, and returns true when
// another pair is available.
func (t *PrefixTrie) Scan() bool {
	ls := len(t.bookmarks)
	if ls == 0 {
		// initialize scan to point to root element
		t.bookmarks = []*pbookmark{&pbookmark{n: t.root}}
		ls++
	}

	// this picks up where it left off
	itop := ls - 1
	bm := t.bookmarks[itop]

outer:
	for {
		// continuing at previous bookmark position
		for ; bm.k < 256; bm.k++ {
			child := bm.n.children[bm.k]
			if child != nil {
				bm = &pbookmark{n: child}
				t.bookmarks = append(t.bookmarks, bm)
				itop++
				if child.valid {
					return true
				}
				continue outer
			}
		}

		for bm.k == 256 {
			if itop--; itop == -1 {
				return false
			}
			bm, t.bookmarks = t.bookmarks[itop], t.bookmarks[:itop+1]
		}

		bm.k++ // next search should start at index after where previous left off
	}

	// never gets here
	return false
}

// Pair returns the key-value pair under the scanning cursor.
func (t *PrefixTrie) Pair() (string, interface{}) {
	var sb strings.Builder
	var n *pnode

	finalIndex := len(t.bookmarks) - 1

	for i, b := range t.bookmarks {
		if i < finalIndex {
			sb.WriteByte(byte(b.k))
			continue
		}
		n = b.n
	}

	return sb.String(), n.value
}

// Text returns the key of the key-value pair under the scanning cursor.
func (t *PrefixTrie) Text() string {
	var sb strings.Builder

	finalIndex := len(t.bookmarks) - 1

	for i, b := range t.bookmarks {
		if i < finalIndex {
			sb.WriteByte(byte(b.k))
			continue
		}
	}

	return sb.String()
}
