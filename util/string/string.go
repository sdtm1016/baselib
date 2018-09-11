package string_util

import (
	"strings"
	"strconv"
	"bytes"
	"math/big"
	"fmt"
	"encoding/base64"
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

func makeAuthHeader(appID, secretKey string) string {
	base64Str := base64.StdEncoding.EncodeToString(
		[]byte(
			fmt.Sprintf("%s:%s", appID, secretKey),
		),
	)
	return fmt.Sprintf("Basic %s", base64Str)
}

func revertValueString(source string) string {
	var target string
	sli := strings.Split(source,"e+")
	val1 := sli[0]
	val2 := sli[1]

	e,err := strconv.Atoi(val2)
	if err!=nil {
		fmt.Println("revertValueString, fail to Atoi(val2), err:", err)
		return ""
	}
	fmt.Println("succeed to Atoi(val2), e:",e)

	sli2 := strings.Split(val1,".")
	length := len(sli2[0]) + e
	tmpS := strings.Replace(val1,".","",-1)
	tmpI,err := strconv.Atoi(tmpS)
	if err!=nil {
		fmt.Println(fmt.Sprintf("revertValueString, fail to Atoi(tmpS), val1:%v, tmpS:%v, err:%v", val1,tmpS,err))
		return ""
	}
	fmt.Println("succeed to Atoi(tmpS), tmpI:",tmpI)

	var offset string
	if len(tmpS) < length{
		gap := length - len(tmpS)
		for i:=0; i<gap; i++{
			offset += "0"
		}
	}

	target = tmpS + offset

	return target
}

func revertTo(source string) string {
	var target string

	tmp := strings.TrimLeft(source,"0")
	target = "0x" + tmp

	return target
}

func revertValueBigIString(source string) string {
	var target string

	length := len(source)
	var gap int
	if length>18{
		gap = length - 18
		sli1 := strings.Split(source,"")
		sli2 := sli1[:gap]
		sli3 := sli1[gap:]
		sli2S := strings.Join(sli2,"")
		sli3S := strings.Join(sli3,"")
		tmp := sli2S + "." + sli3S
		target = strings.TrimRight(tmp,"0")
		target = strings.TrimRight(target,".")
	}else {
		gap = 18 - length
		var tmp  = "0."
		for i:=0; i<gap; i++{
			tmp += "0"
		}
		tmp += source
		target = strings.TrimRight(tmp,"0")
		target = strings.TrimRight(target,".")
	}

	return target
}

func getValueFromInput(){
	//valueTmp := in.Input[length-64:]   ---> bit.Int转化成很大的一个整形数，并将其转化为字符串格式
	valueTmpI,_,err := big.ParseFloat("0000000000000000000000000000000000000000000000005c6d12b6bc1a0000",16,0,0)
	if err!=nil {
		fmt.Println("getValueFromInputT, Fail to ParseFloat(valueTmp), err:",err)
		return
	}

	fmt.Println("getValueFromInputT, valueTmpI:",valueTmpI)
	value := revertValueString(valueTmpI.String())
	fmt.Println("getValueFromInputT, value:",value)
}