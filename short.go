package ab

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"xorm.io/xorm"
)

type SingleModel struct {
	Model             interface{}
	EnablePrivate     bool     // 启用私密模型
	PrivateContextKey string   // 上下文key
	PrivateColName    string   // 字段名
	DisableMethods    []string // get(all) get(single) post put delete
	SearchFields      []string
	Middlewares       []context.Handler
}

type Config struct {
	Party      iris.Party
	StructList []SingleModel
	Engine     *xorm.Engine
}

type modelInfo struct {
	MapName       string
	FullPath      string
	Model         interface{}
	Private       bool
	KeyName       string
	TableColName  string
	StructColName string
	FieldList     TableFieldsResp `json:"field_list"`
	SearchFields  []string
}

type Api struct {
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
