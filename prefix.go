package gotrie

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
	n      *pnode // n points to bookmarked Trie node
	prefix []byte // prefix is the collected key bytes at this node
	k      uint16 // k is part of the bookmark, because it is the next byte to check
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
	indexLastByte := -1
	keyb := []byte(key)

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

// Scan locates the next key-value pair in the Trie. When it finds another pair,
// it returns true; otherwise it returns false.
//
// This works as a continuation, or more specifically as a generator function,
// and only does as much work as required to move the iterator to the next
// key-value pair and return. The first time it is invoked it initializes the
// generator. After it enumerates all key-value pairs in the Trie, it may be
// enumerated again simply by calling this function again.
func (t *PrefixTrie) Scan() bool {
	// As a continuation, this function normally picks back up where it left
	// off. However, if there are no bookmarks, it has either never been
	// executed, or it has already completely enumerated the Trie's contents. In
	// either case, initialize the generator.
	ls := len(t.bookmarks)
	if ls == 0 {
		t.bookmarks = []*pbookmark{&pbookmark{n: t.root}}
		ls++
	}

	itop := ls - 1
	bm := t.bookmarks[itop]

outer:
	for {
		// Look for the next child node from bookmarked node, starting at
		// previous byte from the key.
		for ; bm.k < 256; bm.k++ {
			child := bm.n.children[bm.k]
			if child != nil {
				bm = &pbookmark{
					n:      child,
					prefix: append(bm.prefix, byte(bm.k)),
				}
				t.bookmarks = append(t.bookmarks, bm)
				itop++
				if child.valid {
					return true
				}
				continue outer
			}
		}

		// Current bookmarked node has no additional children, so pop bookmark
		// stack until we find a bookmarked node that has more children to
		// search.
		for bm.k == 256 {
			if itop--; itop == -1 {
				// When the slice of bookmarks is empty, then there are no more
				// key-value pairs to enumerate.
				return false
			}

			// Use top bookmark by popping off the stack of bookmarks.
			bm, t.bookmarks = t.bookmarks[itop], t.bookmarks[:itop+1]
		}

		// The next search must start at the index _after_ the current index.
		bm.k++
	}

	// never gets here
	return false
}

// Pair returns the key-value pair under the scanning cursor.
func (t *PrefixTrie) Pair() (string, interface{}) {
	bm := t.bookmarks[len(t.bookmarks)-1] // top bookmark
	return string(bm.prefix), bm.n.value
}

// Text returns the key of the key-value pair under the scanning cursor.
func (t *PrefixTrie) Text() string {
	return string(t.bookmarks[len(t.bookmarks)-1].prefix)
}
