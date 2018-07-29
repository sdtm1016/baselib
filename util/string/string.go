package string_util

import (
	"strings"
	"strconv"
	"bytes"
)

func IsEmpty(val string) bool{
	return len(strings.TrimSpace(val)) <= 0
}

func RemoveStringsSpace(val string) string{
	// 去除空格  
	val =strings.Replace(val," ","",-1)
	val =strings.Replace(val,"\n","",-1)

	return val
}

func GetString(str interface{}) string {
	if str == nil {
		return ""
	}else if tmp, ok := str.(string); ok {
		return tmp
	}else {
		return ""
	}
}

func ConcatWithSplit(src, dst, split string) string{
	if IsEmpty(dst) {
		return src
	}else if IsEmpty(src){
		return dst
	}else {
		return src + split + dst
	}
}

func DeleteEle(src, split, ele string) string{
	if IsEmpty(src) {
		return src
	}

	srcArr := strings.Split(src, split)
	tmpArr := []string{}
	for _, v := range srcArr {
		if IsEmpty(v) {
			continue
		}
		if v != ele {
			tmpArr = append(tmpArr, v)
		}
	}
	return strings.Join(tmpArr, split)
}

func GetInt(str interface{}) int {
	st := GetString(str)
	i, _ := strconv.Atoi(st)
	return i
}

func SliceToString(slice []string) string{
	if len(slice) <= 0 {
		return ""
	}
	buf := bytes.Buffer{}
	for index, value := range slice{
		buf.WriteString(value)
		if len(slice) > index - 1 {
			buf.WriteString(",")
		}
	}
	return buf.String()
}