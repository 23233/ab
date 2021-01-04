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

func New(c *Config) *RestApi {
	a := new(RestApi)
	a.C = c
	a.checkConfig()
	a.Run()
	return a
}

func (c *RestApi) Run() {
	for _, item := range c.C.StructList {
		model := item.Model
		apiName := c.C.Mdb.TableName(model)
		// 拼接效率高
		p := strings.Join([]string{"/", item.Prefix, apiName, item.Suffix}, "")
		api := c.C.Party.Party(p)

		// resp
		if item.GetAllResponse != nil {
			item.allResp = respItem{
				Has:      true,
				Instance: item.GetAllResponse,
				Fields:   c.tableNameGetNestedStructMaps(reflect.TypeOf(item.GetAllResponse)),
			}
		}
		if item.GetSingleResponse != nil {
			item.singleResp = respItem{
				Has:      true,
				Instance: item.GetSingleResponse,
				Fields:   c.tableNameGetNestedStructMaps(reflect.TypeOf(item.GetSingleResponse)),
			}
		}
		if item.PostResponse != nil {
			item.postResp = respItem{
				Has:      true,
				Instance: item.PostResponse,
				Fields:   c.tableNameGetNestedStructMaps(reflect.TypeOf(item.PostResponse)),
			}
		}
		if item.PutResponse != nil {
			item.putResp = respItem{
				Has:      true,
				Instance: item.PutResponse,
				Fields:   c.tableNameGetNestedStructMaps(reflect.TypeOf(item.PutResponse)),
			}
		}
		if item.DeleteResponse != nil {
			item.deleteResp = respItem{
				Has:      true,
				Instance: item.DeleteResponse,
				Fields:   c.tableNameGetNestedStructMaps(reflect.TypeOf(item.DeleteResponse)),
			}
		}

		item.private = len(item.PrivateContextKey) >= 1 && len(item.PrivateColName) >= 1

		info := modelInfo{
			MapName:   apiName,
			FieldList: c.tableNameReflectFieldsAndTypes(model),
			FullPath:  api.GetRelPath(),
		}
		item.info = info

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
			item.searchFields = result
		}

		// 判断是否还有其他中间件
		if len(item.Middlewares) >= 1 {
			api.Use(item.Middlewares...)
		}

		//
		methods := item.getMethods()

		// 获取全部方法
		if !isContain(methods, "get(all)") {
			var h context.Handler
			if item.GetAllFunc == nil {
				h = c.GetAllFunc
			} else {
				h = item.GetAllFunc
			}
			api.Handle("GET", "/", h)
		}

		// 获取单条
		if !isContain(methods, "get(single)") {
			var h context.Handler
			if item.GetSingleFunc == nil {
				h = c.GetSingle
			} else {
				h = item.GetSingleFunc
			}
			api.Handle("GET", "/{id:uint64}", h)
		}

		// 新增
		if !isContain(methods, "post") {

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
		if !isContain(methods, "put") {
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
		if !isContain(methods, "delete") {

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

func (c *RestApi) pathGetModel(pathName string) SingleModel {
	for _, m := range c.C.StructList {
		if m.info.FullPath == pathName || strings.HasPrefix(pathName, m.info.FullPath) {
			return m
		}
	}
	return SingleModel{}
}

func (c *RestApi) tableNameReflectFieldsAndTypes(table interface{}) tableFieldsResp {
	modelInfo, err := c.C.Mdb.TableInfo(table)
	if err != nil {
		return tableFieldsResp{}
	}
	var resp tableFieldsResp
	// 获取三要素
	values := c.tableNameGetNestedStructMaps(reflect.TypeOf(table))
	resp.Fields = values
	resp.AutoIncrement = modelInfo.AutoIncrement
	resp.Version = modelInfo.Version
	resp.Deleted = modelInfo.Deleted
	resp.Created = modelInfo.Created
	resp.Updated = modelInfo.Updated
	return resp
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
			d.MapName = c.C.Mdb.GetColumnMapper().Obj2Table(field.Name)
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
		d.MapName = c.C.Mdb.GetColumnMapper().Obj2Table(field.Name)
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
	for _, item := range c.C.StructList {
		if item.info.MapName == tableName {
			return item.Model, nil
		}
	}
	return nil, errors.New("未找到模型")
}

// 通过模型名获取模型信息
func (c *RestApi) tableNameGetModelInfo(tableName string) (SingleModel, error) {
	for _, l := range c.C.StructList {
		if l.info.MapName == tableName {
			return l, nil
		}
	}
	return SingleModel{}, errors.New("未找到模型")
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

	fields := cb.info.FieldList

	for _, column := range fields.Fields {
		if column.MapName != fields.AutoIncrement {
			if column.MapName == fields.Updated || column.MapName == fields.Deleted {
				continue
			}
			if len(fields.Created) >= 1 {
				var equal = false
				for k := range fields.Created {
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
