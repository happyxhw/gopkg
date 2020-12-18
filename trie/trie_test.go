package trie

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"testing"
)

func TestTrie_Search(t1 *testing.T) {
	trie := NewTrie()
	words := []string{"中国人民", "中国", "人民", "伟大"}
	for _, item := range words {
		trie.Insert(item, "1")
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

func loadData() []string {
	var res []string
	file, err := os.Open("./x.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = append(res, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return res
}

func loadSearchData() []string {
	var res []string
	file, err := os.Open("./y.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = append(res, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return res
}

type record struct {
	SearchWord string
	RecWord    string
	TrieWord   string
}

func TestPlumWork(t1 *testing.T) {
	res := loadData()
	trie := NewTrie()
	for _, item := range res {
		if item != "" {
			word := strings.Split(item, ",")
			if len(word) == 2 {
				trie.Insert(word[0], word[1])
			}
		}
	}

	searchRecord := loadSearchData()
	var cleanRecord []*record

	for _, item := range searchRecord {
		if item != "" {
			word := strings.Split(item, ",")
			if s := strings.TrimSpace(word[0]); s != "" {
				t := record{
					SearchWord: s,
					RecWord:    strings.TrimSpace(word[1]),
				}
				cleanRecord = append(cleanRecord, &t)
			}
		}
	}

	sortMap := map[string]int{
		"suffix":   1,
		"series":   2,
		"brand":    3,
		"category": 4,
		"tag":      5,
	}

	for i, item := range cleanRecord {
		words := strings.Split(item.SearchWord, " ")
		var searchWords, searchEnWords []string
		filterMap := make(map[string]bool)
		for _, w := range words {
			w := strings.TrimSpace(w)
			if isEng(w) {
				searchEnWords = append(searchEnWords, w)
			} else {
				searchWords = append(searchWords, w)
			}
		}
		var resList []*SearchResult
		if len(searchWords) > 0 {
			resList = trie.Search(strings.Join(searchWords, " "), SuffixMaxSearch)
		}
		for _, w := range searchEnWords {
			ret, wordType := trie.Find(w)
			if ret {
				resList = append(resList, &SearchResult{
					WordType: wordType,
					Word:     w,
				})
			}
		}
		result := make([]string, 0, len(resList))
		// suffix series brand category tag

		sort.Slice(resList, func(i, j int) bool {
			iRank, jRank := sortMap[resList[i].WordType], sortMap[resList[j].WordType]
			return iRank < jRank
		})

		for _, item2 := range resList {
			if _, ok := filterMap[item2.Word]; !ok {
				result = append(result, item2.Word+"-"+item2.WordType)
				if item2.WordType == "suffix" && strings.HasSuffix(item.SearchWord, item2.Word) {
					if len(result) > 1 {
						t := result[0]
						result[0] = result[len(result)-1]
						result[len(result)-1] = t
					}
				}
				filterMap[item2.Word] = true
			}
		}
		if len(result) > 0 {
			cleanRecord[i].TrieWord = strings.Join(result, ", ")
		}
	}
	f, err := os.Create("test.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, _ = f.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(f)
	for _, item := range cleanRecord {
		s := []string{
			item.SearchWord, item.TrieWord, item.RecWord,
		}
		_ = w.Write(s)
	}
	w.Flush()

}

func isEng(word string) bool {
	for _, item := range strings.ToLower(word) {
		if item < 'a' || item > 'z' {
			return false
		}
	}
	return true
}

func TestPlumWork1(t1 *testing.T) {
	res := loadData()
	trie := NewTrie()
	for _, item := range res {
		if item != "" {
			word := strings.Split(item, ",")
			if len(word) == 2 {
				trie.Insert(word[0], word[1])
			}
		}

	}
	//
	content := "2.55 chanel 腰包"
	resList := trie.Search(content, SuffixMaxSearch)
	for _, item := range resList {
		fmt.Printf("%+v\n", item)
	}
}
