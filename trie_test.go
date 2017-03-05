package LevenshteinTrie

import (
	"bufio"
	"io"
	"os"
	"testing"
)

var tree *TrieNode

func getfile() (io.ReadCloser, error) {
	filename := "./w1_fixed.txt"
	file, err := os.Open(filename)
	return file, err
}

func TestInsert(t *testing.T) {
	tree = NewTrie()
	file, err := getfile()
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if err != nil {
			break
		}

		tree.InsertText(line)
	}
}

func TestPrefixSearch(t *testing.T) {
	expected := []string{
		"zygodactyl",
		"zygoma",
		"zygomatic",
		"zygomaticus",
		"zygomycetes",
		"zygon",
		"zygosity",
		"zygote",
		"zygote-specific",
		"zygotes",
		"zygourakis",
	}
	query := "zygo"
	words := tree.Suffix(query)
	for _, e := range expected {
		if !contains(words, e) {
			t.Errorf("Missing word: %s", e)
		} else {
			t.Logf("Found: %s\n", e)
		}
	}
}

func TestExactSearch(t *testing.T) {
	query := "zygon"
	node := tree.Get(query)
	if node.GetText() != query {
		t.Error("Node not found")
	}

	if tree.Get("zygo") != nil {
		t.Error("Should be nil object")
	}
}

func contains(words []string, word string) bool {
	for _, w := range words {
		if w == word {
			return true
		}
	}
	return false
}
func TestLevenshteinSearch(t *testing.T) {
	expected := []struct {
		query    string
		result   []QueryResult
		distance int
	}{
		{"accidens", []QueryResult{
			{"accidens", 0, nil},
			{"accident", 1, nil},
		}, 1},
	}
	for _, e := range expected {
		results := tree.Levenshtein(e.query, e.distance)
		if !containsq(results, e.query) {
			t.Errorf("Missing term: %s", e.query)
		} else {
			t.Logf("Looking for: %s, Found: %s\n", e.query, results)
		}
	}
}

func containsq(results []QueryResult, word string) bool {
	for _, r := range results {
		if r.Val == word {
			return true
		}
	}
	return false
}

func TestTreeWithMeta(t *testing.T) {
	r := NewTrie()
	r.Add("romane", 1)
	r.Add("romanus", 2)
	r.Add("romulus", 3)
	r.Add("ruber", 4)
	r.Add("rubens", 5)
	r.Add("rubicon", 6)
	r.Add("rubicundus", 7)

	meta := r.Get("ruber").GetInfo()
	if meta.(int) != 4 {
		t.Error("Metadata not match")
	}

	qs := r.Levenshtein("rubycon", 1)
	meta = qs[0].Node.GetInfo()

	if meta.(int) != 6 {
		t.Error("Metadata not match")
	}

	r.Add("ruber", 10)
	meta = r.Get("ruber").GetInfo()
	if meta.(int) != 10 {
		t.Error("Metadata not match")
	}
}
