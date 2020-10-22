package ab

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"xorm.io/xorm"
)

// 错误返回
func fastError(err error, ctx iris.Context, msg ...string) {
	ctx.StatusCode(iris.StatusBadRequest)
	var m string
	if err == nil {
		m = "请求解析出错"
	} else {
		m = err.Error()
	}
	if len(msg) >= 1 {
		m = msg[0]
	}
	_, _ = ctx.JSON(iris.Map{
		"detail": m,
	})
	return
}

// 获取所有 分页 页码用page标识
func GetAllFunc(ctx iris.Context) {
	page := ctx.URLParamIntDefault("page", 1)
	if page > 100 {
		page = 100
	}
	pageSize := ctx.URLParamIntDefault("page_size", 20)
	if pageSize > 100 {
		pageSize = 100
	}
	model := nowApi.pathGetModel(ctx.Path())

	// 解析出order by
	descField := ctx.URLParam("order_desc")
	orderBy := ctx.URLParam("order")
	// 从url中解析出filter
	filterList := filterMatch(ctx.URLParams(), model.FieldList.Fields)

	privateName := ctx.Values().Get(model.KeyName)
	start := (page - 1) * pageSize
	end := page * (pageSize * 2)

	var base = func() *xorm.Session {
		var d *xorm.Session
		d = nowApi.Config.Engine.Table(model.MapName)
		if model.Private {
			d = d.Where(fmt.Sprintf("%s = ?", model.TableColName), privateName)
		}
		if len(orderBy) >= 1 {
			d = d.OrderBy(orderBy)
		} else if len(descField) >= 1 {
			d = d.Desc(descField)
		}
		return d
	}

	where := func() *xorm.Session {
		var d *xorm.Session
		d = base()
		if len(model.FieldList.Deleted) >= 1 {
			d = base().Where(fmt.Sprintf("%s = ? OR %s IS NULL", model.FieldList.Deleted, model.FieldList.Deleted), "0001-01-01 00:00:00")
		}
		if len(filterList) >= 1 {
			for k, v := range filterList {
				d = d.Where(fmt.Sprintf("%s = ?", k), v)
			}
		}
		return d
	}

	// 获取总数量
	allCount, err := where().Count()
	if err != nil {
		fastError(err, ctx)
		return
	}

	// 获取内容
	dataList := make([]map[string]string, 0)
	if allCount >= 1 {
		if len(model.FieldList.AutoIncrement) >= 1 && len(filterList) < 1 && len(orderBy) < 1 && len(descField) < 1 {
			dataList, err = where().And(fmt.Sprintf("%s between ? and ?", model.FieldList.AutoIncrement), start, end).Limit(pageSize).QueryString()
		} else {
			dataList, err = where().Limit(pageSize, start).QueryString()
		}
		if err != nil {
			fastError(err, ctx)
			return
		}
	}
	result := iris.Map{
		"page_size": pageSize,
		"page":      page,
		"all":       allCount,
		"data":      dataList,
	}
	if len(descField) >= 1 {
		result["desc_field"] = descField
	}
	if len(orderBy) >= 1 {
		result["order"] = orderBy
	}
	if len(filterList) >= 1 {
		result["filter"] = filterList
	}
	_, _ = ctx.JSON(result)

}

// 单个 /{id:uint64}
func GetSingle(ctx iris.Context) {
	id, err := ctx.Params().GetUint64("id")
	if err != nil {
		fastError(err, ctx, "参数错误")
		return
	}
	model := nowApi.pathGetModel(ctx.Path())
	privateName := ctx.Values().Get(model.KeyName)
	newData := nowApi.newModel(model.MapName)

	var base = func() *xorm.Session {
		if model.Private {
			return nowApi.Config.Engine.Table(newData).Where(fmt.Sprintf("%s = ?", model.TableColName), privateName)
		}
		return nowApi.Config.Engine.Table(newData)
	}
	has, err := base().ID(id).Get(newData)
	if err != nil || has == false {
		fastError(err, ctx, "未找到数据")
		return
	}
	_, _ = ctx.JSON(newData)
}

// 新增数据
func AddData(ctx iris.Context) {
	model := nowApi.pathGetModel(ctx.Path())
	newInstance, err := nowApi.getCtxValues(model.MapName, ctx)
	if err != nil {
		fastError(err, ctx)
		return
	}
	if model.Private {
		privateName := ctx.Values().Get(model.KeyName)
		private := newInstance.Elem().FieldByName(model.StructColName)
		switch private.Type().String() {
		case "string":
			private.SetString(privateName.(string))
			break
		case "int", "int8", "int16", "int32", "int64", "time.Duration":
			private.SetInt(int64(privateName.(int)))
			break
		case "uint", "uint8", "uint16", "uint32", "uint64":
			private.SetUint(uint64(privateName.(int)))
			break
		default:
			fastError(err, ctx, "私密参数解析错误")
			return
		}
	}

	singleData := newInstance.Interface()

	aff, err := nowApi.Config.Engine.Table(model.MapName).InsertOne(singleData)
	if err != nil || aff == 0 {
		fastError(err, ctx, "新增数据失败")
		return
	}
	_, _ = ctx.JSON(iris.Map{})
}

// 编辑数据 /{id:uint64}
func EditData(ctx iris.Context) {
	model := nowApi.pathGetModel(ctx.Path())
	privateName := ctx.Values().Get(model.KeyName)
	id, err := ctx.Params().GetUint64("id")
	if err != nil {
		fastError(err, ctx, "参数错误")
		return
	}

	var base = func() *xorm.Session {
		if model.Private {
			return nowApi.Config.Engine.Table(model.MapName).Where(fmt.Sprintf("%s = ?", model.TableColName), privateName)
		}
		return nowApi.Config.Engine.Table(model.MapName)
	}
	// 先获取数据是否存在
	has, err := base().Where("id = ?", id).Exist()
	if err != nil {
		fastError(err, ctx)
		return
	}
	if has != true {
		fastError(err, ctx, "获取数据失败")
		return
	}
	newInstance, err := nowApi.getCtxValues(model.MapName, ctx)
	if err != nil {
		fastError(err, ctx)
		return
	}

	if model.Private {
		privateName := ctx.Values().Get(model.KeyName)
		private := newInstance.Elem().FieldByName(model.StructColName)
		switch private.Type().String() {
		case "string":
			private.SetString(privateName.(string))
			break
		case "int", "int8", "int16", "int32", "int64", "time.Duration":
			private.SetInt(int64(privateName.(int)))
			break
		case "uint", "uint8", "uint16", "uint32", "uint64":
			private.SetUint(uint64(privateName.(int)))
			break
		default:
			fastError(err, ctx, "私密参数解析错误")
			return
		}
	}

	// 全量更新
	singleData := newInstance.Interface()
	aff, err := nowApi.Config.Engine.Table(model.MapName).ID(id).AllCols().Update(singleData)
	if err != nil || aff < 1 {
		fastError(err, ctx, "更新数据失败")
		return
	}
	_, _ = ctx.JSON(iris.Map{})
}

// 删除数据 /{id:uint64}
func DeleteData(ctx iris.Context) {
	// 先获取
	model := nowApi.pathGetModel(ctx.Path())
	privateName := ctx.Values().Get(model.KeyName)
	id, err := ctx.Params().GetUint64("id")
	newData := nowApi.newModel(model.MapName)

	if err != nil {
		fastError(err, ctx, "参数错误")
		return
	}
	var base = func() *xorm.Session {
		if model.Private {
			return nowApi.Config.Engine.Table(newData).Where(fmt.Sprintf("%s = ?", model.TableColName), privateName)
		}
		return nowApi.Config.Engine.Table(newData)
	}
	// 先获取数据是否存在

	has, err := base().ID(id).Get(newData)
	if err != nil {
		fastError(err, ctx)
		return
	}
	if has != true {
		fastError(err, ctx, "获取数据失败")
		return
	}
	// 进行删除
	aff, err := base().ID(id).Delete(newData)
	if err != nil || aff < 1 {
		fastError(err, ctx, "删除数据失败")
		return
	}
	_, _ = ctx.JSON(iris.Map{})

}
