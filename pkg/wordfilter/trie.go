package wordfilter

// TrieTreeFilter 字典树过滤器
type TrieTreeFilter struct {
	root *TrieNode
}

type TrieNode struct {
	Leaf bool
	Next map[rune]*TrieNode
	Rune rune
}

func (n *TrieNode) AddChild(c rune) *TrieNode {
	if n.Next == nil {
		n.Leaf = false
		n.Next = make(map[rune]*TrieNode)
	}
	if next, ok := n.Next[c]; ok {
		return next
	}
	next := &TrieNode{Leaf: true, Rune: c}
	n.Next[c] = next
	return next
}

func (n *TrieNode) FindChild(c rune) *TrieNode {
	if n.Leaf {
		return nil
	}
	if node, ok := n.Next[c]; ok {
		return node
	}
	return nil
}

func (f *TrieTreeFilter) AddWord(word string) {
	node := f.root
	chars := []rune(word)
	for _, c := range chars {
		node = node.AddChild(c)
	}
}

func (f *TrieTreeFilter) Build(dict []string) {
	for _, word := range dict {
		f.AddWord(word)
	}
}

func (f *TrieTreeFilter) DoFilter(text string) ([]string, string) {
	sensitives := make([]string, 0)
	runes := []rune(text)
	length := len(runes)
	for i := 0; i < length; i++ {
		// 用句子中的每一个字符作为单词的开始位置，去root匹配
		child := f.root.FindChild(runes[i])
		// 第一个字符不匹配，跳过
		if child == nil {
			continue
		}
		var j int
		// 从该字符开始继续匹配
		for j = i + 1; j < length && child != nil; j++ {
			// 已经匹配到了叶子节点，表示成功匹配了整个单词
			if child.Leaf {
				sensitives = append(sensitives, string(runes[i:j]))
				replaceRunes(runes, '*', i, j)
				// i 到 j 的字符已经替换过了,须注意的是此时的j已经在下一个字符了，所以需要-1
				i = j - 1
				break
			}
			child = child.FindChild(runes[j])
		}
		// 字符串末尾的特殊情况
		if j == length && child != nil && child.Leaf {
			sensitives = append(sensitives, string(runes[i:j]))
			replaceRunes(runes, '*', i, j)
			i = j
		}
	}
	return sensitives, string(runes)
}

func replaceRunes(runes []rune, replacer rune, start, end int) {
	for i := start; i < end; i++ {
		runes[i] = replacer
	}
}
