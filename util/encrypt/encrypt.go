package encrypt

import (
	"crypto/sha1"
	"fmt"
	"crypto/md5"
	"io"
	"encoding/base64"
	"time"
)

const (
	TIME_BASE_TIME_FMT = "2006-01-02 15:04:05"
)

 /**
 产生散列值 sha1.New()，sha1.Write(bytes)，然后sha1.Sum([]byte{})
  */
func GetSha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	//%x 散列结果格式化为 16 进制字符串
	return fmt.Sprintf("%x", bs)
}

func GetMd5(s string) string {
	w := md5.New()
	io.WriteString(w, s)

	return fmt.Sprintf("%x", w.Sum(nil))
}

func ConvertBase64(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

func TimeFormat(time time.Time) string {
	return time.Format(TIME_BASE_TIME_FMT)
}

func GetTimestamp() int64{
	return time.Now().UnixNano() / int64(time.Millisecond)
}
