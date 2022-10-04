package backstage

import (
	"time"
)

func StartHeartbeat() {
	go func() {
		for {
			// 假设定时任务是
			// 从redis取出注册的ip:port
			// online和offline都要发送请求
			time.Sleep(5)
		}
	}()
}
