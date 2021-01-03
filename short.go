package ab

import (
	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"time"
	"xorm.io/xorm"
)

type SingleModel struct {
	Middlewares        []context.Handler
	Prefix             string                 // 路由前缀
	Suffix             string                 // 路由后缀
	Model              interface{}            // xorm model
	model              modelInfo              // model info
	PrivateContextKey  string                 // 上下文key
	PrivateColName     string                 // 数据库字段名
	ExtraFilters       map[string]interface{} //
	EnableMethods      []string               //
	DisableMethods     []string               // get(all) get(single) post put delete
	AllowSearchFields  []string               // 搜索的字段 struct名称
	GetAllFunc         func(ctx iris.Context) // 覆盖获取全部方法
	GetAllResponse     interface{}            // 获取所有返回的内容替换 仅替换data数组 同名替换
	GetSingleFunc      func(ctx iris.Context) // 覆盖获取单条方法
	GetSingleResponse  interface{}            // 获取单个返回的内容替换
	PostFunc           func(ctx iris.Context) // 覆盖新增方法
	PostValidator      interface{}            // 新增自定义验证器
	PostResponse       interface{}            // 新增返回内容
	PutFunc            func(ctx iris.Context) // 覆盖修改方法
	PutValidator       interface{}            // 修改验证器
	PutResponse        interface{}            // 修改返回内容
	DeleteFunc         func(ctx iris.Context) // 覆盖删除方法
	DeleteValidator    interface{}            // 删除验证器
	DeleteResponse     interface{}            // 删除返回内容
	CacheTime          time.Time              //
	GetAllCacheTime    time.Time              //
	GetSingleCacheTime time.Time              //
}

type MysqlConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
	PoolSize int
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	Db       int
	PoolSize int
}

type IrisInstance struct {
	Party iris.Party
	App   iris.Application
}

type MysqlInstance struct {
	MysqlConfig
	Mdb *xorm.Engine
}

type RedisInstance struct {
	RedisConfig
	Rdb *redis.Client
}

type Config struct {
	IrisInstance
	MysqlInstance
	RedisInstance
	StructList []SingleModel
}

type respItem struct {
	Has    bool
	Model  interface{}
	Fields []structInfo
}

type modelInfo struct {
	MapName   string
	FullPath  string
	FieldList TableFieldsResp
}

type RestApi struct {
	Config     *Config
	ModelLists []modelInfo
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

type TableFieldsResp struct {
	Fields        []structInfo    `json:"fields"`
	AutoIncrement string          `json:"autoincr"`
	Updated       string          `json:"updated"`
	Deleted       string          `json:"deleted"`
	Created       map[string]bool `json:"created"`
	Version       string          `json:"version"`
}
