package trie

const (
	SuffixSearch = iota + 1
	SuffixMaxSearch
	SuffixMinSearch
	AllSearch
)

type Node struct {
	childrenMap map[rune]*Node
	isWordEnd   bool
	wordType    string
}

type SearchResult struct {
	Start    int
	End      int
	WordType string
	Word     string
}

// Trie tree
type Trie struct {
	root *Node
}

// NewTrie init trie tree
func NewTrie() *Trie {
	return &Trie{
		root: &Node{
			childrenMap: make(map[rune]*Node),
		},
	}
}

// Insert insert word to the tree
func (t *Trie) Insert(word string, wordType string) {
	current := t.root
	for _, item := range word {
		if _, ok := current.childrenMap[item]; !ok {
			current.childrenMap[item] = &Node{
				childrenMap: make(map[rune]*Node),
			}
		}
		current = current.childrenMap[item]
	}
	current.isWordEnd = true
	current.wordType = wordType
}

// Find certain word
func (t *Trie) Find(word string) (bool, string) {
	current := t.root
	for _, item := range word {
		if _, ok := current.childrenMap[item]; !ok {
			return false, ""
		}
		current = current.childrenMap[item]
	}
	return current.isWordEnd, current.wordType
}

// Search tree of content
// 词典： “中国人民”，“人民”，“中国”
// 搜索词：“中国人民”
// SuffixSearch     : 中国，中国人民
// SuffixMaxSearch  : 中国人民
// SuffixMinSearch  : 中国，人民
// AllSearch        : 中国，人民，中国人民
func (t *Trie) Search(content string, searchType int) []*SearchResult {
	var res []*SearchResult
	s := []rune(content)
	n := len(s)
	i := 0
LOOP:
	for i < n {
		gap := 0
		cur, start := i, i
		current := t.root
		var temp []*SearchResult
		for {
			if current == nil || cur >= n || current.childrenMap[s[cur]] == nil {
				if searchType == SuffixMaxSearch && len(temp) > 0 {
					res = append(res, temp[len(temp)-1])
				}
				i++
				goto LOOP
			} else {
				item := current.childrenMap[s[cur]]
				gap++
				current = item
				cur++
				if item.isWordEnd {
					end := start + gap
					v := SearchResult{
						Start:    start,
						End:      end,
						WordType: item.wordType,
						Word:     string(s[start:end]),
					}
					switch searchType {
					case SuffixSearch:
						res = append(res, &v)
						i = end
					case SuffixMaxSearch:
						temp = append(temp, &v)
						i = end
						if current == nil || i>=n || cur >= n || current.childrenMap[s[cur]] == nil {
							if len(temp) > 0 {
								res = append(res, temp[len(temp)-1])
							}
							goto LOOP
						}
					case SuffixMinSearch:
						res = append(res, &v)
						i = end
						goto LOOP
					case AllSearch:
						res = append(res, &v)
						if current == nil || cur >= n || current.childrenMap[s[cur]] == nil {
							i++
							goto LOOP
						}
					}
				}
			}
		}

	}
	return res
}
