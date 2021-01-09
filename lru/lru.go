package lru

import "sync"

// DLinkedNode double link list
type DLinkedNode struct {
	Key   int
	Value int
	Pre   *DLinkedNode
	Post  *DLinkedNode
}

// LRU cache
type LRU struct {
	sync.Mutex

	Cache    map[int]*DLinkedNode
	Head     *DLinkedNode
	Tail     *DLinkedNode
	Count    int
	Capacity int
}

// NewLRU init lru
func NewLRU(capacity int) *LRU {
	lru := LRU{
		Cache:    make(map[int]*DLinkedNode),
		Head:     &DLinkedNode{},
		Tail:     &DLinkedNode{},
		Capacity: capacity,
	}
	lru.Head.Post = lru.Tail
	lru.Tail.Pre = lru.Head
	return &lru
}

func (lru *LRU) Get(key int) int {
	lru.Lock()
	defer lru.Unlock()
	if node, ok := lru.Cache[key]; ok {
		lru.moveToFirst(node)
		return node.Value
	} else {
		return -1
	}
}

func (lru *LRU) Put(key int, value int) {
	lru.Lock()
	defer lru.Unlock()
	if node, ok := lru.Cache[key]; ok {
		node.Value = value
		lru.moveToFirst(node)
	} else {
		node := &DLinkedNode{
			Key:   key,
			Value: value,
		}

		lru.Cache[key] = node
		lru.addNode(node)
		lru.Count++
		if lru.Count > lru.Capacity {
			last := lru.popLast()
			delete(lru.Cache, last.Key)
			lru.Count--
		}
	}
}

func (lru *LRU) addNode(node *DLinkedNode) {
	node.Pre = lru.Head
	node.Post = lru.Head.Post
	lru.Head.Post.Pre = node
	lru.Head.Post = node
}

func (lru *LRU) delNode(node *DLinkedNode) {
	pre := node.Pre
	post := node.Post

	pre.Post = post
	post.Pre = pre

	node.Pre = nil
	node.Post = nil
}

func (lru *LRU) moveToFirst(node *DLinkedNode) {
	lru.delNode(node)
	lru.addNode(node)
}

func (lru *LRU) popLast() *DLinkedNode {
	last := lru.Tail.Pre
	lru.delNode(last)
	return last
}
