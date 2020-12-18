package ab

import "github.com/kataras/iris/v12"

// 全局访问中间件 优先级最高
type GlobalPreMiddlewareProcess interface {
	ApiGlobalPreMiddleware(ctx iris.Context)
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
