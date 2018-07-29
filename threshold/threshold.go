package threshold

import (
	"time"
	"sync"
	"errors"
	"fmt"
	"sync/atomic"
)


/**
1.发请求前需要先获取令牌
2.限定某时间段内的发放的令牌数量
3.任务执行完成后，使用定时器不断重置令牌
4.如果当前goroutine数量过多时也不重置令牌
 */

//限流
type threshold struct {
	D      time.Duration //周期是D，
	C      int64         //限制一个周期最多操作C次
	Mu     sync.Mutex
	Token  chan bool //令牌池
	num    int64     //当前的goroutine数量
	maxNum int64     //允许工作goroutine最大数量
}

//如果两个周期后还没有申请到令牌，就报错超时
//目前用不到，如果限制routine最大数量需要靠这来监控
var (
	ErrApplyTimeout = errors.New("apply token time out")
)

func NewThrottle(D time.Duration, C, maxNum int64) *threshold {
	instance := &threshold{
		D:      D,
		C:      C,
		Token:  make(chan bool, C),
		maxNum: maxNum,
	}
	go instance.reset()
	return instance
}

//每周期重新填充一次令牌池
func (t *threshold) reset() {
	ticker := time.NewTicker(t.D)
	for _ = range ticker.C {
		//goroutine数量不超过最大数量时再填充令牌池
		if t.num >= t.maxNum {
			continue
		}
		t.Mu.Lock()
		supply := t.C - int64(len(t.Token))
		fmt.Printf("reset token:%d\n", supply)
		for supply > 0 {
			t.Token <- true
			supply--
		}
		t.Mu.Unlock()
	}
}

//申请令牌，如果过两个周期还没申请到就报超时退出
func (t *threshold) ApplyToken() (bool, error) {
	select {
	case <-t.Token:
		return true, nil
	case <-time.After(t.D * 2):
		return false, ErrApplyTimeout
	}
}

func (t *threshold) Work(job func()) {
	if ok, err := t.ApplyToken(); !ok {
		fmt.Println(err)
		return
	}
	go func() {
		atomic.AddInt64(&t.num, 1)
		defer atomic.AddInt64(&t.num, -1)
		job()
	}()
}