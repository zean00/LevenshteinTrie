LevenshteinTrie
===============

Change : Add support for metadata / information on inserted node

###Additional API
```
//Insert text with metadata
t.Add("word",10)
t.Add("another",struct{}{})

//Get node from exact word 
node := t.Get("word")
//Get metadata from node
info := node.GetInfo()
text := node.GetText()

//QueryResult with node
results := t.Levenshtein("crave", 1)
node := results[0].Node

//Get nodes from suffix search
var nodes []*TrieNode
nodes := t.NodesSuffix("cratyl")
for _, n := range nodes {
	fmt.Println(n.GetText())
}
```

A Trie data structure that allows for fuzzy string matching

This is the Go version of a python program written by Steve Hanov in his [blog post](http://stevehanov.ca/blog/index.php?id=114)

This is finished, but not tested.
[![Build Status](https://drone.io/github.com/jamra/LevenshteinTrie/status.png)](https://drone.io/github.com/jamra/LevenshteinTrie/latest)
###How it works

 - It is a basic [Trie](http://en.wikipedia.org/wiki/Trie).

 - You can search for all words that are suffixes of a string. 

 - You can also search for words within a certain edit distance of a string. The algorithm memoizes the Levenshtein algorithm when it recursively iterates through the Trie nodes. This speeds up the Levenshtein matches hugely.

###Example

```
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	lt "github.com/jamra/LevenshteinTrie"
)

func main() {
	file, err := os.Open("./w1_fixed.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	t := lt.NewTrie()
	for scanner.Scan() {
		line := scanner.Text()
		if err != nil {
			break
		}

		t.InsertText(line)
	} 
	results := t.Suffix("cratyl")
	fmt.Println(results)
	results2 := t.Levenshtein("crave", 1)
	fmt.Println(results2)
}
```
