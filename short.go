package ab

import (
	_ctx "context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/pkg/errors"
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
	ExtraFilters       map[string]interface{} //
	AllowMethods       []string               //
	DisableMethods     []string               // get(all) get(single) post put delete
	AllowSearchFields  []string               // 搜索的字段 struct名称
	searchFields       []string               //
	GetAllFunc         func(ctx iris.Context) // 覆盖获取全部方法
	GetAllResponse     interface{}            // 获取所有返回的内容替换 仅替换data数组 同名替换
	allResp            respItem               //
	GetSingleFunc      func(ctx iris.Context) // 覆盖获取单条方法
	GetSingleResponse  interface{}            // 获取单个返回的内容替换
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
	CacheTime          time.Time              //
	GetAllCacheTime    time.Time              //
	GetSingleCacheTime time.Time              //
}

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
	result := make([]string, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}
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
