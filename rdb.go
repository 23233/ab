package ab

import (
	"context"
	"github.com/OneOfOne/xxhash"
	"github.com/jxskiss/base62"
	"io"
	"strconv"
	"strings"
	"time"
)

// 此文件主要放redis相关操作

// genRedisKey 生成redis存储的key 尽量的短 所以使用 xxhash后进行base62
// 生成key需要的参数为 所有请求参数与额外参数 额外参数可以为用户id等
func genRedisKey(ReqParams string, otherInfo ...string) string {
	d := make([]string, 1+len(otherInfo))
	d = append(d, ReqParams)
	d = append(d, otherInfo...)
	origin := strings.Join(d, "")
	// 进行hash
	h := xxhash.New64()
	r := strings.NewReader(origin)
	_, _ = io.Copy(h, r)
	keyInt := h.Sum64()
	// 进行base62
	return base62.EncodeToString([]byte(strconv.FormatUint(keyInt, 10)))
}

// saveToRedis 响应体保存到redis当中
func (c *RestApi) saveToRedis(ctx context.Context, keyName string, data string, expireTime time.Duration) error {
	return c.C.Rdb.Set(ctx, keyName, data, expireTime).Err()
}

// 删除key
func (c *RestApi) deleteToRedis(ctx context.Context, keyName string) error {
	return c.C.Rdb.Del(ctx, keyName).Err()
}
