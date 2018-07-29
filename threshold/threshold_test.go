package threshold

import (
	"time"
	"fmt"
	"testing"
)

func TestNewThrottle(t *testing.T) {
	ts := NewThrottle(time.Second, 10, 20) //每秒10次，同时最多20个routine存在
	for {
		ts.Work(doWork)
	}
}

//真正的工作函数 假设每个需要执行5秒
func doWork() {
	fmt.Println(time.Now())
	<-time.After(5 * time.Second)
}
