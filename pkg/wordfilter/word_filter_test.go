package wordfilter

import "testing"

func TestTrieTreeFilterEnglish(t *testing.T) {
	dict := []string{"stupid", "moron", "fuck", "bitch"}
	filter := NewTrieTreeFilter()
	filter.Build(dict)
	t.Log(filter.root.Next)
	text1 := "you stupid fucking moron"
	text2 := "you mother fucker"
	text3 := "bitchfucker"

	sensitive1, replaced1 := filter.DoFilter(text1)
	if replaced1 != "you ****** ****ing *****" {
		t.Fail()
	}
	t.Log(sensitive1)
	t.Log(replaced1)

	sensitive2, replaced2 := filter.DoFilter(text2)
	if replaced2 != "you mother ****er" {
		t.Fail()
	}
	t.Log(sensitive2)
	t.Log(replaced2)

	sensitive3, replaced3 := filter.DoFilter(text3)
	if replaced3 != "*********er" {
		t.Fail()
	}
	t.Log(sensitive3)
	t.Log(replaced3)
}

func TestTrieTreeFilterChinese(t *testing.T) {
	dict := []string{"傻逼", "你妈的", "他妈的", "SB"}
	filter := NewTrieTreeFilter()
	filter.Build(dict)
	t.Log(filter.root.Next)
	text1 := "你他妈的就是个傻逼"
	text2 := "你妈的你就是个傻逼"
	text3 := "他妈的SB"

	sensitive1, replaced1 := filter.DoFilter(text1)
	if replaced1 != "你***就是个**" {
		t.Fail()
	}
	t.Log(sensitive1)
	t.Log(replaced1)

	sensitive2, replaced2 := filter.DoFilter(text2)
	if replaced2 != "***你就是个**" {
		t.Fail()
	}
	t.Log(sensitive2)
	t.Log(replaced2)

	sensitive3, replaced3 := filter.DoFilter(text3)
	if replaced3 != "*****" {
		t.Fail()
	}
	t.Log(sensitive3)
	t.Log(replaced3)
}
