package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/karrick/gotrie"
)

func main() {
	// build a new Trie from standard input lines
	t := gotrie.NewPrefixTrie()
	t.Insert("", 0)

	scanner := bufio.NewScanner(os.Stdin)
	var line int
	for scanner.Scan() {
		line++
		t.Insert(scanner.Text(), line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	// Enumerate through Trie in sorted order
	for t.Scan() {
		fmt.Println(t.Text())
	}

	deleteFromTrie(t, "")

	deleteFromTrie(t, "romane")
	deleteFromTrie(t, "romanus")
	deleteFromTrie(t, "romulus")

	deleteFromTrie(t, "romane")
	deleteFromTrie(t, "romanus")
	deleteFromTrie(t, "romulus")
	deleteFromTrie(t, "rubens")
	deleteFromTrie(t, "ruber")
	deleteFromTrie(t, "rubicon")
	deleteFromTrie(t, "rubicundus")

	deleteFromTrie(t, "")
}

func deleteFromTrie(t *gotrie.PrefixTrie, key string) {
	fmt.Fprintf(os.Stderr, "key: %q; ok: %t\n", key, t.Delete(key))
}
