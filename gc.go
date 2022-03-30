package cache

import "time"

type garcoll struct {
	interval time.Duration // 自动清理周期
	stop     chan bool     // 停止回收的通道
}

func (gc *garcoll) Run(c *Cache) {
	ticker := time.NewTicker(gc.interval)
	for {
		select {
		case <-ticker.C:
			c.delete()
		case <-gc.stop:
			ticker.Stop()
			return
		}
	}
}

// func (j *janitor) Run(c *Cache) {
// 	ticker := time.NewTicker(j.interval)
// 	for {
// 		select {
// 		case <-ticker.C:
// 			c.DeleteExpired()
// 		case <-j.stop:
// 			ticker.Stop()
// 			return
// 		}
// 	}
// }

// func runJanitor(c *Cache, d time.Duration) {
// 	j := &janitor{
// 		interval: d,
// 		stop:     make(chan bool),
// 	}
// 	c.janitor = j
// 	go j.Run(c)
// }
