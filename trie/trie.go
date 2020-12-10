package trie

type Node struct {
	childrenMap map[rune]*Node
	isWordEnd   bool
	wordType    int
}

type SearchResult struct {
	Start    int
	End      int
	WordType int
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
func (t *Trie) Insert(word string, wordType int) {
	s := []rune(word)
	current := t.root
	for _, item := range s {
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
func (t *Trie) Find(word string) bool {
	s := []rune(word)
	current := t.root
	for _, item := range s {
		if _, ok := current.childrenMap[item]; !ok {
			return false
		}
		current = current.childrenMap[item]
	}
	if current.isWordEnd {
		return true
	}
	return false
}

// Search tree of content
// isSuffix: 词典 “中国人民”，“人民”，“中国”
// true: “中国人民” 将优先匹配 “中国” “中国人民” 不匹配 “人民”
// false: “中国人民” 将匹配 “中国” “中国人民” 和 “人民”
func (t *Trie) Search(content string, isSuffix bool) []*SearchResult {
	var res []*SearchResult
	s := []rune(content)
	n := len(s)
	i := 0
	current := t.root
	for i < n {
		gap := 1
		cur, start := i, i
		for current != nil && cur < n {
			if item, ok := current.childrenMap[s[cur]]; !ok {
				i++
				current = t.root
				break
			} else {
				gap++
				current = item
				cur++
				if item.isWordEnd {
					end := start + gap - 1
					v := SearchResult{
						Start:    start,
						End:      end,
						WordType: item.wordType,
						Word:     string(s[start:end]),
					}
					res = append(res, &v)
					if isSuffix {
						i = end
					}
				}
			}
		}

	}
	return res
}
