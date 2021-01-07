package ab

import "log"

func (c *RestApi) checkConfig() {
	c.C.MysqlInstance.check()
	hasCache := false
	for _, model := range c.C.Models {
		if model.getAllListCacheTime() >= 1 || model.getSingleCacheTime() >= 1 {
			hasCache = true
			break
		}
	}
	if hasCache {
		c.C.RedisInstance.check()
	}
	if c.C.ErrorTrace == nil {
		c.C.ErrorTrace = func(err error, event, from, router string) {
			log.Printf("[ab][%s] error:%s event:%s from:%s ", router, err, event, from)
		}
	}
}
