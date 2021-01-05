package ab

import (
	_ctx "context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/pkg/errors"
	"strings"
	"time"
	"xorm.io/xorm"
)

type SingleModel struct {
	Middlewares        []context.Handler
	Prefix             string                 // 路由前缀
	Suffix             string                 // 路由后缀
	Model              interface{}            // xorm model
	info               modelInfo              //
	private            bool                   //
	PrivateContextKey  string                 // 上下文key
	PrivateColName     string                 // 数据库字段名
	AllowMethods       []string               //
	DisableMethods     []string               // get(all) get(single) post put delete
	AllowSearchFields  []string               // 搜索的字段 struct名称
	searchFields       []string               //
	GetAllFunc         func(ctx iris.Context) // 覆盖获取全部方法
	GetAllResponse     interface{}            // 获取所有返回的内容替换 仅替换data数组 同名替换
	GetAllExtraFilters map[string]string      // 额外的固定过滤 key(数据库列名) 和 value 若与请求过滤重复则覆盖 优先级最高
	allResp            respItem               //
	GetSingleFunc      func(ctx iris.Context) // 覆盖获取单条方法
	GetSingleResponse  interface{}            // 获取单个返回的内容替换
	SingleExtraFilters map[string]string      // 额外的固定过滤 key(数据库列名) 和 value 若与请求过滤重复则覆盖 优先级最高
	singleResp         respItem               //
	PostFunc           func(ctx iris.Context) // 覆盖新增方法
	PostValidator      interface{}            // 新增自定义验证器
	PostResponse       interface{}            // 新增返回内容
	postResp           respItem               //
	PutFunc            func(ctx iris.Context) // 覆盖修改方法
	PutValidator       interface{}            // 修改验证器
	PutResponse        interface{}            // 修改返回内容
	putResp            respItem               //
	DeleteFunc         func(ctx iris.Context) // 覆盖删除方法
	DeleteValidator    interface{}            // 删除验证器
	DeleteResponse     interface{}            // 删除返回内容
	deleteResp         respItem               //
	CacheTime          time.Duration          //
	GetAllCacheTime    time.Duration          //
	GetSingleCacheTime time.Duration          //
	DelayDeleteTime    time.Duration          // 延迟多久双删
	MaxPageSize        int                    //
	MaxPageCount       int                    //
}

// getMethods 初始化请求方法 返回数组
func (c *SingleModel) getMethods() []string {
	if len(c.AllowMethods) >= 1 {
		return c.AllowMethods
	}
	m := c.initMethods()
	if len(c.DisableMethods) >= 1 {
		for _, method := range c.DisableMethods {
			if _, ok := m[method]; ok {
				delete(m, method)
				continue
			}
		}
	}
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// initMethods 初始化请求方法 返回map
func (c *SingleModel) initMethods() map[string]string {
	// get(all) get(single) post put delete
	return map[string]string{
		"get(all)":    "get(all)",
		"get(single)": "get(single)",
		"post":        "post",
		"put":         "put",
		"delete":      "delete",
	}
}

// getPage 获取最大限制的页码和每页数量
func (c *SingleModel) getPage() (int, int) {
	maxPageCount := c.MaxPageCount
	if maxPageCount < 1 {
		maxPageCount = 100
	}
	maxPageSize := c.MaxPageSize
	if maxPageSize < 1 {
		maxPageSize = 100
	}
	return maxPageCount, maxPageSize
}

// getDelayDeleteTime 获取延迟删除时间
func (c *SingleModel) getDelayDeleteTime() time.Duration {
	if c.DelayDeleteTime >= 1 {
		return c.DelayDeleteTime
	}
	return 500 * time.Millisecond
}

// getAllListCacheTime 获取列表缓存时间
func (c *SingleModel) getAllListCacheTime() time.Duration {
	if c.GetAllCacheTime >= 1 {
		return c.GetAllCacheTime
	}
	return c.CacheTime
}

// getSingleCacheTime 获取单条缓存时间
func (c *SingleModel) getSingleCacheTime() time.Duration {
	if c.GetSingleCacheTime >= 1 {
		return c.GetSingleCacheTime
	}
	return c.CacheTime
}

// getAllExtraParams 额外参数解析成url形式
func (c *SingleModel) getAllExtraParams() string {
	var s strings.Builder
	for k, v := range c.GetAllExtraFilters {
		s.WriteString(fmt.Sprintf("%s=%s", k, v))
	}
	return s.String()
}

// getSingleExtraParams 额外参数解析成url形式
func (c *SingleModel) getSingleExtraParams() string {
	var s strings.Builder
	for k, v := range c.SingleExtraFilters {
		s.WriteString(fmt.Sprintf("%s=%s", k, v))
	}
	return s.String()
}

type MysqlConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
	PoolSize int
	ShowSql  bool
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	Db       int
	PoolSize int
}

type MysqlInstance struct {
	MysqlConfig
	Mdb *xorm.Engine
}

func (c *MysqlInstance) check() {
	if c.Mdb == nil {
		if len(c.MysqlConfig.Host) < 1 {
			panic("[mysql] config mysql config or engine instance must be need")
		} else {
			c.connect()
		}
	}
}
func (c *MysqlInstance) connect() {
	// database 连接器
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4", c.Username, c.Password, c.Host, c.Port, c.DbName)
	engine, err := xorm.NewEngine("mysql", dbUrl)
	if err != nil {
		panic(err)
	}
	engine.SetMaxOpenConns(c.PoolSize)
	engine.ShowSQL(c.ShowSql)

	c.Mdb = engine
	err = c.ping()
	if err != nil {
		panic(errors.Wrap(err, "[mysql] ping fail"))
	}
}
func (c *MysqlInstance) ping() error {
	return c.Mdb.Ping()
}

type RedisInstance struct {
	RedisConfig
	Rdb *redis.Client
}

func (c *RedisInstance) check() {
	if c.Rdb == nil {
		if len(c.RedisConfig.Host) < 1 {
			panic("[redis]config or instance must needs")
		} else {
			c.connect()
		}
	}
}
func (c *RedisInstance) connect() {
	poolSize := c.PoolSize
	if poolSize < 1 {
		poolSize = 100
	}
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: c.Password,
		DB:       c.Db,
		PoolSize: poolSize,
	})
	c.Rdb = client
	err := c.ping()
	if err != nil {
		panic(errors.Wrap(err, "[redis] ping fail"))
	}
}
func (c *RedisInstance) ping() error {
	var ctx = _ctx.Background()
	return c.Rdb.Ping(ctx).Err()
}

// config
type Config struct {
	Party iris.Party
	MysqlInstance
	RedisInstance
	StructList []SingleModel
}

type modelInfo struct {
	MapName   string
	FullPath  string
	FieldList tableFieldsResp
}

type RestApi struct {
	C *Config
}

// 模型信息
type structInfo struct {
	Name         string `json:"name"`
	Types        string `json:"types"`
	MapName      string `json:"map_name"`
	XormTags     string `json:"xorm_tags"`
	ValidateTags string `json:"validate_tags"`
	CommentTags  string `json:"comment_tags"`
	AttrTags     string `json:"attr_tags"`
}

type tableFieldsResp struct {
	Fields        []structInfo    `json:"fields"`
	AutoIncrement string          `json:"autoincr"`
	Updated       string          `json:"updated"`
	Deleted       string          `json:"deleted"`
	Created       map[string]bool `json:"created"`
	Version       string          `json:"version"`
}

type respItem struct {
	Has      bool
	Instance interface{}
	Fields   []structInfo
}
