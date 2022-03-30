package cache

type trie struct {
	children [60]*trie
	isEnd    bool
}

func newTrie() *trie {
	return &trie{}
}

func (t *trie) insert(key string) {
	node := t
	for _, ch := range key {
		ch -= 'A'
		if node.children[ch] == nil {
			node.children[ch] = &trie{}
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

func (t *trie) searchPrefix(prefix string) *trie {
	node := t
	for _, ch := range prefix {
		ch -= 'A'
		if node.children[ch] == nil {
			return nil
		}
		node = node.children[ch]
	}
	return node
}

func (t *trie) startsWithPrefix(prefix string) bool {
	return t.searchPrefix(prefix) != nil
}
