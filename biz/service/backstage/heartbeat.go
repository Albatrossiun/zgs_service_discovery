package backstage

import (
	"time"
)

func StartHeartbeat() {
	go func() {
		for {
			time.Sleep(5)
		}
	}()
}
