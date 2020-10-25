package ab

import "github.com/kataras/iris/v12"
import "github.com/kataras/iris/v12/context"

// 全局访问中间件 优先级最高
type GlobalPreMiddlewareProcess interface {
	ApiGlobalPreMiddleware(ctx iris.Context)
}

// 中间件列表 优先级第二
type MiddlewareListProcess interface {
	ApiMiddleware() []context.Handler
}

// 私密访问
type PrivateAccessProcess interface {
	// 上下文获取私密条件内容的key
	ApiPrivateContextKey() string
	// 数据列名
	ApiPrivateTableColName() string
}

// 禁止方法生成 默认生成 get(all) get(single) post put delete
type DisableMethodsProcess interface {
	ApiDisableMethods() []string
}

// 搜索字段
type SearchFieldsProcess interface {
	ApiSearchFields() []string
}

// 方法单独的中间件
type GetAllPreMiddlewareProcess interface {
	ApiGetAllPreMiddleware(ctx iris.Context)
}
type GetSinglePreMiddlewareProcess interface {
	ApiGetSinglePreMiddleware(ctx iris.Context)
}
type PostPreMiddlewareProcess interface {
	ApiPostPreMiddleware(ctx iris.Context)
}
type PutPreMiddlewareProcess interface {
	ApiPutPreMiddleware(ctx iris.Context)
}
type PutDeleteMiddlewareProcess interface {
	ApiDeletePreMiddleware(ctx iris.Context)
}

// 覆盖方法
type GetAllProcess interface {
	ApiGetAll(ctx iris.Context)
}
type GetSingleProcess interface {
	ApiGetSingle(ctx iris.Context)
}
type PostProcess interface {
	ApiPost(ctx iris.Context)
}
type PutProcess interface {
	ApiPut(ctx iris.Context)
}
type DeleteProcess interface {
	ApiDelete(ctx iris.Context)
}
