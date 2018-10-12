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
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		t.Insert(scanner.Text(), struct{}{})
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	// Enumerate through Trie in sorted order
	for t.Scan() {
		fmt.Println(t.Text())
	}
}
