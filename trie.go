package LevenshteinTrie

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

func min(a ...int) int {
	min := int(^uint(0) >> 1) // largest int
	for _, i := range a {
		if i < min {
			min = i
		}
	}
	return min
}

func max(a ...int) int {
	max := int(0)
	for _, i := range a {
		if i > max {
			max = i
		}
	}
	return max
}

//TrieNode trie structure
type TrieNode struct {
	letter   rune //Equivalent to int32
	children map[rune]*TrieNode
	final    bool
	text     string
	meta     interface{}
}

func (t *TrieNode) String() string {
	s := fmt.Sprintf("%U\n", t.letter)
	for _, v := range t.children {
		s += fmt.Sprintf("-%#v\n", v)
	}
	return s
}

//NewTrie create new Levenshtein trie
func NewTrie() *TrieNode {
	return &TrieNode{children: make(map[rune]*TrieNode)}
}

//Add add node with metadata
func (t *TrieNode) Add(text string, meta interface{}) {
	if t == nil {
		return
	}
	text = strings.ToLower(text)
	currNode := t //Starts at root
	for i, w := 0, 0; i < len(text); i += w {
		runeValue, width := utf8.DecodeRuneInString(text[i:])
		final := false
		if width+i == len(text) {
			final = true
		}
		w = width

		currNode = newTrieNode(currNode, runeValue, final, text, meta)
	}
}

//GetInfo get node metadata
func (t *TrieNode) GetInfo() interface{} {
	return t.meta
}

//GetText get text string
func (t *TrieNode) GetText() string {
	return t.text
}

//InsertText insert text as node
func (t *TrieNode) InsertText(text string) {
	if t == nil {
		return
	}
	text = strings.ToLower(text)
	currNode := t //Starts at root
	for i, w := 0, 0; i < len(text); i += w {
		runeValue, width := utf8.DecodeRuneInString(text[i:])
		final := false
		if width+i == len(text) {
			final = true
		}
		w = width

		currNode = newTrieNode(currNode, runeValue, final, text, nil)
	}
}

//Get node with exact string
func (t *TrieNode) Get(text string) *TrieNode {
	ns := t.NodesSuffix(text)
	if len(ns) == 0 {
		return nil
	}

	for _, v := range ns {
		if v.text == text {
			return v
		}
	}
	return nil
}

//Suffix query on suffix
func (t *TrieNode) Suffix(query string) []string {
	ns := t.NodesSuffix(query)
	result := make([]string, len(ns))
	for i, v := range ns {
		result[i] = v.text
	}
	return result
}

//NodesSuffix query node on suffix
func (t *TrieNode) NodesSuffix(query string) []*TrieNode {
	var curr *TrieNode
	var ok bool

	curr = t
	//first, find the end of the prefix
	for _, letter := range query {
		if curr != nil {
			curr, ok = curr.children[letter]
			if ok {
				//do nothing
			}

		} else {
			return nil
		}
	}

	candidates := getsuffixr(curr)

	return candidates
}

func newTrieNode(t *TrieNode, runeValue rune, final bool, text string, meta interface{}) *TrieNode {
	node, exists := t.children[runeValue]
	if meta == nil {
		meta = text
	}
	if !exists {
		node = &TrieNode{letter: runeValue, children: make(map[rune]*TrieNode), meta: meta}
		t.children[runeValue] = node
	} else {
		node.meta = meta
	}
	if final {
		node.final = true
		node.text = text
	}
	return node
}

func getsuffixr(n *TrieNode) []*TrieNode {
	if n == nil {
		return nil
	}

	var candidates []*TrieNode

	if n.final == true {
		candidates = append(candidates, n)
	}

	for _, childNode := range n.children {
		candidates = append(candidates, getsuffixr(childNode)...)
	}
	return candidates
}

//QueryResult query result
type QueryResult struct {
	Val      string
	Distance int
	Node     *TrieNode
}

func (q QueryResult) String() string {
	return fmt.Sprintf("Val: %s, Dist: %d\n", q.Val, q.Distance)
}

type byDistance []QueryResult

func (a byDistance) Len() int           { return len(a) }
func (a byDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDistance) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

//Levenshtein query with Levenshtein distance
func (t *TrieNode) Levenshtein(text string, distance int) []QueryResult {

	//initialize the first row for the dynamic programming alg
	l := utf8.RuneCount([]byte(text))
	currentRow := make([]int, l+1)

	for i := 0; i < len(currentRow); i++ {
		currentRow[i] = i
	}

	var candidates []QueryResult

	for letter, childNode := range t.children {
		candidates = append(candidates, searchlevr(childNode, currentRow, letter, []rune(text), distance)...)
	}

	sort.Sort(byDistance(candidates))
	return candidates
}

func searchlevr(n *TrieNode, prevRow []int, letter rune, text []rune, maxDistance int) []QueryResult {
	columns := len(prevRow)
	currentRow := make([]int, columns)

	currentRow[0] = prevRow[0] + 1

	for col := 1; col < columns; col++ {
		if text[col-1] == letter {
			currentRow[col] = prevRow[col-1]
			continue
		}
		insertCost := currentRow[col-1] + 1
		deleteCost := prevRow[col] + 1
		replaceCost := prevRow[col-1] + 1

		currentRow[col] = min(insertCost, deleteCost, replaceCost)
	}

	var candidates []QueryResult

	distance := currentRow[len(currentRow)-1]
	if distance <= maxDistance && n.final == true {
		candidates = append(candidates, QueryResult{Val: n.text, Distance: distance, Node: n})
	}
	mi := min(currentRow[1:]...)
	if mi <= maxDistance {
		for l, childNode := range n.children {
			candidates = append(candidates, searchlevr(childNode, currentRow, l, text, maxDistance)...)
		}
	}
	return candidates
}
