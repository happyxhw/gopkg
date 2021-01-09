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
func (t *Trie) Insert(word, wordType string) {
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
func (t *Trie) Find(word string) (isWordEnd bool, wordType string) {
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
						if current == nil || i >= n || cur >= n || current.childrenMap[s[cur]] == nil {
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

func MinEditDistance(word1, word2 string) int {
	m, n := len(word1), len(word2)
	cur, pre := make([]int, n+1), make([]int, n+1)
	for j := 1; j < n+1; j++ {
		pre[j] = j
	}
	for i := 1; i <= m; i++ {
		cur[0] = i
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				cur[j] = pre[j-1]
			} else {
				cur[j] = min(pre[j-1], min(pre[j], cur[j-1])) + 1
			}
		}
		for i := range pre {
			pre[i] = 0
		}
		pre, cur = cur, pre
	}

	return pre[n]
}

func CommonPrefix(word1, word2 string) int {
	m, n := len(word1), len(word2)
	var c int
	for i := 0; i < m && i < n; i++ {
		if word1[i] == word2[i] {
			c++
		} else {
			break
		}
	}
	return c
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func VagueSearch(keyword string, enWords []string) string {
	if len(enWords) == 0 {
		return ""
	}
	// 前缀匹配
	n := float32(len(enWords))
	var word string
	var maxScore float32

	for _, item := range enWords {
		c := float32(CommonPrefix(keyword, item))
		score := c / n
		// 匹配的前缀占一半比例
		if score > 0.5 && c/float32(len(item)) > 0.5 {
			if score > maxScore {
				maxScore = score
				word = item
			}
		}
	}
	if word != "" {
		return word
	}

	// 编辑距离
	minDis := 1000
	for _, item := range enWords {
		dis := MinEditDistance(keyword, item)
		if dis == 0 {
			return item
		}
		if dis <= 2 && dis < minDis {
			minDis = dis
			word = item
		}
	}
	return word
}
