package string_util

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestDeleteEle(t *testing.T) {
	str1 := "1,2,3,4"
	assert.Equal(t, DeleteEle(str1, ",", "1"), "2,3,4")
}

func TestConcatWithSplit(t *testing.T) {
	var(
		src = "guo"
		dst = "bin"
		split = "@"
	)
	ret := ConcatWithSplit(src,dst,split)
	fmt.Printf("ret:%s",ret)
}

func TestRemoveStringsSpace(t *testing.T) {
	str := "xhayzbhh jjn kkoq   ggbb\n"
	ret := RemoveStringsSpace(str)
	fmt.Printf("ret:%s",ret)
}

func TestSliceToString(t *testing.T) {
	//sli := []string{"hello","guo","bin"}
	var sli = make([]string,0)
	sli = append(sli,"hello")
	sli = append(sli,"guo")
	sli = append(sli,"bin")
	ret := SliceToString(sli)
	fmt.Printf("ret:%s",ret)
}