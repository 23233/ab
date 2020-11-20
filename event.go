package ab

import "github.com/kataras/iris/v12"

// 全局访问中间件 优先级最高
type GlobalPreMiddlewareProcess interface {
	ApiGlobalPreMiddleware(ctx iris.Context)
}

// 私密访问
type PrivateAccessProcess interface {
	// 返回上下文key 数据列名
	ApiPrivate() (string, string)
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
