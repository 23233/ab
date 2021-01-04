package ab

func (c *RestApi) checkConfig() {
	c.C.MysqlInstance.check()
	c.C.RedisInstance.check()
}
