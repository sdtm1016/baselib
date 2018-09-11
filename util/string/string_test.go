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

func TestBigIntToString(t *testing.T) {
	ret1 := makeAuthHeader("3308060dfd8ff","9060b7de4bab94e95a1758021a73d920")
	ret2 := revertValueString("1.02e+20")
	ret3 := revertTo("00000000000000000000000012587c9ece7ccedaddda06e92f8659ffb68efaac")
	ret4 := revertValueBigIString("12000000000000000000")//12
	ret5 := revertValueBigIString("12030000000000000001")//12.030000000000000001
	ret6 := revertValueBigIString("12030000000000000")//0.01203
	ret7 := revertValueBigIString("6660000000000000000")//6.66

	fmt.Println("-------------------ret1:",ret1)
	fmt.Println("-------------------ret2:",ret2)
	fmt.Println("-------------------ret3:",ret3)
	fmt.Println("-------------------ret4:",ret4)
	fmt.Println("-------------------ret5:",ret5)
	fmt.Println("-------------------ret6:",ret6)
	fmt.Println("-------------------ret7:",ret7)

	fmt.Println("----------------------")

	getValueFromInput()
}