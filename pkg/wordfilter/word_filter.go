package wordfilter

type Filter interface {
	// Build 用词典构建敏感词过滤器
	Build(dict []string)

	// DoFilter 执行敏感词过滤，返回匹配的敏感词和替换后的字符串
	DoFilter(text string) ([]string, string)
}

func NewTrieTreeFilter() *TrieTreeFilter {
	return &TrieTreeFilter{root: &TrieNode{Leaf: true}}
}
