package trie

import (
	"fmt"
	"testing"
)

func TestTrie_Search(t1 *testing.T) {
	trie := NewTrie()
	words := []string{"中国人民", "中国", "人民", "伟大"}
	for _, item := range words {
		trie.Insert(item, 1)
	}
	content := "中国人民真伟大"

	fmt.Println("SuffixSearch")
	resList := trie.Search(content, SuffixSearch)
	for _, item := range resList {
		fmt.Printf("%+v\n", item)
	}
	fmt.Println("SuffixMaxSearch")
	resList = trie.Search(content, SuffixMaxSearch)
	for _, item := range resList {
		fmt.Printf("%+v\n", item)
	}
	fmt.Println("SuffixMinSearch")
	resList = trie.Search(content, SuffixMinSearch)
	for _, item := range resList {
		fmt.Printf("%+v\n", item)
	}
	fmt.Println("AllSearch")
	resList = trie.Search(content, AllSearch)
	for _, item := range resList {
		fmt.Printf("%+v\n", item)
	}

}
