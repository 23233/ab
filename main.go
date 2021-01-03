package ab

import (
	"github.com/23233/sv"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func New(c Config) *RestApi {
	a := new(RestApi)
	a.Config = &c
	a.Run()
	return a
}

func (c *RestApi) Run() {
	for _, item := range c.Config.StructList {
		model := item.Model
		apiName := c.Config.Engine.TableName(model)
		api := c.Config.Party.Party("/" + apiName)

		info := modelInfo{
			MapName:       apiName,
			Model:         model,
			Private:       len(item.PrivateContextKey) >= 1 && len(item.PrivateColName) >= 1,
			KeyName:       item.PrivateContextKey,
			StructColName: item.PrivateColName,
			FieldList:     c.tableNameReflectFieldsAndTypes(apiName),
			FullPath:      api.GetRelPath(),
		}
		if item.GetAllResponse != nil {
			info.GetAllResp = respItem{
				Has:    true,
				Model:  item.GetAllResponse,
				Fields: c.tableNameGetNestedStructMaps(reflect.TypeOf(item.GetAllResponse)),
			}
		}
		if item.GetSingleResponse != nil {
			info.GetSingleResp = respItem{
				Has:    true,
				Model:  item.GetSingleResponse,
				Fields: c.tableNameGetNestedStructMaps(reflect.TypeOf(item.GetSingleResponse)),
			}
		}
		if item.PostResponse != nil {
			info.PostResp = respItem{
				Has:    true,
				Model:  item.PostResponse,
				Fields: c.tableNameGetNestedStructMaps(reflect.TypeOf(item.PostResponse)),
			}
		}
		if item.PutResponse != nil {
			info.PutResp = respItem{
				Has:    true,
				Model:  item.PutResponse,
				Fields: c.tableNameGetNestedStructMaps(reflect.TypeOf(item.PutResponse)),
			}
		}
		if item.DeleteResponse != nil {
			info.DeleteResp = respItem{
				Has:    true,
				Model:  item.DeleteResponse,
				Fields: c.tableNameGetNestedStructMaps(reflect.TypeOf(item.DeleteResponse)),
			}
		}

		if len(item.AllowSearchFields) >= 1 {
			var result []string
			for _, f := range item.AllowSearchFields {
				for _, field := range info.FieldList.Fields {
					if field.Name == f || field.MapName == f {
						result = append(result, field.MapName)
						break
					}
				}
			}
			info.SearchFields = result
		}

		if info.Private {
			for _, field := range info.FieldList.Fields {
				if field.Name == item.PrivateColName {
					info.TableColName = field.MapName
					break
				}
			}
		}

		c.ModelLists = append(c.ModelLists, info)

		// 判断使用拥有前置访问中间件
		if processor, ok := model.(GlobalPreMiddlewareProcess); ok {
			api.Use(processor.ApiGlobalPreMiddleware)
		}

		// 判断是否还有其他中间件
		if len(item.Middlewares) >= 1 {
			api.Use(item.Middlewares...)
		}

		// 获取全部方法
		if !isContain(item.DisableMethods, "get(all)") {
			var h context.Handler
			if item.GetAllFunc == nil {
				h = c.GetAllFunc
			} else {
				h = item.GetAllFunc
			}
			api.Handle("GET", "/", h)
		}

		// 获取单条
		if !isContain(item.DisableMethods, "get(single)") {
			var h context.Handler
			if item.GetSingleFunc == nil {
				h = c.GetSingle
			} else {
				h = item.GetSingleFunc
			}
			api.Handle("GET", "/{id:uint64}", h)
		}

		// 新增
		if !isContain(item.DisableMethods, "post") {

			var h context.Handler
			if item.PostFunc == nil {
				h = c.AddData
			} else {
				h = item.PostFunc
			}
			route := api.Handle("POST", "/", h)

			// 判断是否有自定义验证器
			if item.PostValidator != nil {
				route.Use(sv.Run(item.PostValidator))
			}
		}

		// 修改
		if !isContain(item.DisableMethods, "put") {
			var h context.Handler
			if item.PutFunc == nil {
				h = c.EditData
			} else {
				h = item.PutFunc
			}
			route := api.Handle("PUT", "/{id:uint64}", h)
			// 判断是否有自定义验证器
			if item.PutValidator != nil {
				route.Use(sv.Run(item.PutValidator))
			}
		}

		// 删除
		if !isContain(item.DisableMethods, "delete") {

			var h context.Handler
			if item.DeleteFunc == nil {
				h = c.DeleteData
			} else {
				h = item.DeleteFunc
			}
			route := api.Handle("DELETE", "/{id:uint64}", h)
			// 判断是否有自定义验证器
			if item.DeleteValidator != nil {
				route.Use(sv.Run(item.DeleteValidator))
			}
		}

	}

}

func (c *RestApi) pathGetModel(pathName string) modelInfo {
	for _, m := range c.ModelLists {
		if m.FullPath == pathName || strings.HasPrefix(pathName, m.FullPath) {
			return m
		}
	}
	return modelInfo{}
}

func (c *RestApi) tableNameReflectFieldsAndTypes(tableName string) TableFieldsResp {
	for _, item := range c.Config.StructList {
		model := item.Model
		routerName := c.Config.Engine.TableName(model)
		if routerName == tableName {
			modelInfo, err := c.Config.Engine.TableInfo(model)
			if err != nil {
				return TableFieldsResp{}
			}
			var resp TableFieldsResp
			// 获取三要素
			values := c.tableNameGetNestedStructMaps(reflect.TypeOf(model))
			resp.Fields = values
			resp.AutoIncrement = modelInfo.AutoIncrement
			resp.Version = modelInfo.Version
			resp.Deleted = modelInfo.Deleted
			resp.Created = modelInfo.Created
			resp.Updated = modelInfo.Updated
			return resp
		}
	}
	return TableFieldsResp{}

}

// 通过模型名获取所有列信息 名称 类型 xorm tag validator comment
func (c *RestApi) tableNameGetNestedStructMaps(r reflect.Type) []structInfo {
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	if r.Kind() != reflect.Struct {
		return nil
	}
	v := reflect.New(r).Elem()
	result := make([]structInfo, 0)
	for i := 0; i < r.NumField(); i++ {
		field := r.Field(i)
		v := reflect.Indirect(v).FieldByName(field.Name)
		fieldValue := v.Interface()
		var d structInfo

		switch fieldValue.(type) {
		case time.Time, time.Duration:
			d.Name = field.Name
			d.Types = field.Type.String()
			d.XormTags = field.Tag.Get("xorm")
			d.ValidateTags = field.Tag.Get("validate")
			d.CommentTags = field.Tag.Get("comment")
			d.AttrTags = field.Tag.Get("attr")
			d.MapName = c.Config.Engine.GetColumnMapper().Obj2Table(field.Name)
			result = append(result, d)
			continue
		}
		if field.Type.Kind() == reflect.Struct {
			values := c.tableNameGetNestedStructMaps(field.Type)
			result = append(result, values...)
			continue
		}
		d.Name = field.Name
		d.Types = field.Type.String()
		d.MapName = c.Config.Engine.GetColumnMapper().Obj2Table(field.Name)
		d.XormTags = field.Tag.Get("xorm")
		d.CommentTags = field.Tag.Get("comment")
		d.AttrTags = field.Tag.Get("attr")
		d.ValidateTags = field.Tag.Get("validate")
		result = append(result, d)
	}
	return result
}

// 通过模型名获取实例
func (c *RestApi) tableNameGetModel(tableName string) (interface{}, error) {
	for _, item := range c.ModelLists {
		if item.MapName == tableName {
			return item, nil
		}
	}
	return nil, errors.New("未找到模型")
}

// 通过模型名获取模型信息
func (c *RestApi) tableNameGetModelInfo(tableName string) (modelInfo, error) {
	for _, l := range c.ModelLists {
		if l.MapName == tableName {
			return l, nil
		}
	}
	return modelInfo{}, errors.New("未找到模型")
}

// 获取内容
func (c *RestApi) getValue(ctx iris.Context, k string) string {
	var b string
	b = ctx.PostValueTrim(k)
	if len(b) < 1 {
		b = ctx.FormValue(k)
	}
	b = strings.Trim(b, " ")
	return b
}

// 对应关系获取
func (c *RestApi) getCtxValues(routerName string, ctx iris.Context) (reflect.Value, error) {
	// 先获取到字段信息
	cb, err := c.tableNameGetModelInfo(routerName)
	if err != nil {
		return reflect.Value{}, err
	}
	t := reflect.TypeOf(cb.Model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	newInstance := reflect.New(t)

	for _, column := range cb.FieldList.Fields {
		if column.MapName != cb.FieldList.AutoIncrement {
			if column.MapName == cb.FieldList.Updated || column.MapName == cb.FieldList.Deleted {
				continue
			}
			if len(cb.FieldList.Created) >= 1 {
				var equal = false
				for k := range cb.FieldList.Created {
					if column.MapName == k {
						equal = true
						break
					}
				}
				if equal {
					continue
				}
			}
			content := c.getValue(ctx, column.MapName)
			if len(content) < 1 {
				continue
			}
			switch column.Types {
			case "string":
				newInstance.Elem().FieldByName(column.Name).SetString(content)
				continue
			case "int", "int8", "int16", "int32", "int64", "time.Duration":
				d, err := strconv.ParseInt(content, 10, 64)
				if err != nil {
					log.Printf("解析出int出错")
				}
				newInstance.Elem().FieldByName(column.Name).SetInt(d)
				continue
			case "uint", "uint8", "uint16", "uint32", "uint64":
				d, err := strconv.ParseUint(content, 10, 64)
				if err != nil {
					log.Println("解析出uint出错")
				}
				newInstance.Elem().FieldByName(column.Name).SetUint(d)
				continue
			case "float32", "float64":
				d, err := strconv.ParseFloat(content, 64)
				if err != nil {
					log.Println("解析出float出错")
				}
				newInstance.Elem().FieldByName(column.Name).SetFloat(d)
				continue
			case "bool":
				d, err := parseBool(content)
				if err != nil {
					log.Println("解析出bool出错")
				}
				newInstance.Elem().FieldByName(column.Name).SetBool(d)
				continue
			case "time", "time.Time":
				var tt reflect.Value
				// 判断是否是字符串
				if IsNum(content) {
					// 这里需要转换成时间
					d, err := strconv.ParseInt(content, 10, 64)
					if err != nil {
						return reflect.Value{}, errors.Wrap(err, "time change to int error")
					}
					tt = reflect.ValueOf(time.Unix(d, 0))
				} else {
					formatTime, err := time.ParseInLocation("2006-01-02 15:04:05", content, time.Local)
					if err != nil {
						return reflect.Value{}, errors.Wrap(err, "time parse location error")
					}
					tt = reflect.ValueOf(formatTime)
				}
				newInstance.Elem().FieldByName(column.Name).Set(tt)
				continue
			}
		}
	}

	return newInstance, nil
}

// 模型反射一个新
func (c *RestApi) newModel(routerName string) interface{} {
	cb, _ := c.tableNameGetModelInfo(routerName)
	return c.newType(cb.Model)
}

// 反射一个新数据
func (c *RestApi) newType(input interface{}) interface{} {
	t := reflect.TypeOf(input)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	newInstance := reflect.New(t)
	return newInstance.Interface()
}
