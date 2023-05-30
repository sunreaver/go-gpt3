package gpt3

import (
	"math/rand"
	"time"
)

// RetryHandle 重试机制的处理类型
type RetryHandle func() error

// 尝试重试机制
// fn：处理事件方法
// retries：重试次数
// delay：重试间隔
// 返回值为可选的 err 或 nil
func retry(fn RetryHandle, retries int, sleep time.Duration) error {
	if err := fn(); err != nil {
		retries--
		if retries <= 0 {
			return err
		}

		sleep += (time.Duration(rand.Int63n(int64(sleep)))) / 2
		time.Sleep(sleep)

		return retry(fn, retries, 2*sleep)
	}

	return nil
}
