package trie

import (
	"fmt"
	"testing"
)

func TestTrie_Search(t1 *testing.T) {
	trie := NewTrie()
	words := []string{"中国人民", "中华", "傻逼", "中国", "你真逗比", "国人", "人民"}
	for _, item := range words {
		trie.Insert(item, 1)
	}
	fmt.Println(trie.Find("中国人"))
	content := "中国人shi傻逼笑死我了，哈哈，中华，xda你真逗比sdfas中国人民"
	resList := trie.Search(content, false)
	for _, item := range resList {
		fmt.Printf("%+v\n", item)
	}
}
