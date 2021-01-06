package ab

func (c *RestApi) checkConfig() {
	c.C.MysqlInstance.check()
	hasCache := false
	for _, model := range c.C.StructList {
		if model.getAllListCacheTime() >= 1 || model.getSingleCacheTime() >= 1 {
			hasCache = true
			break
		}
	}
	if hasCache {
		c.C.RedisInstance.check()
	}
}
