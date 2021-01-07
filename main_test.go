package ab

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/iris-contrib/httpexpect/v2"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/httptest"
	_ "github.com/mattn/go-sqlite3"
	"testing"
	"time"
	"xorm.io/xorm"
)

type testModel struct {
	Id   uint64 `xorm:"autoincr pk unique" json:"id"`
	Name string `xorm:"varchar(10)" json:"name"`
	Age  uint64 `json:"age"`
	Desc string `xorm:"varchar(20)" json:"desc"`
	Code uint64 `xorm:"" json:"code"`
}

func TestNew(t *testing.T) {
	app := iris.New()
	prefix := "/api/v1"

	p := app.Party(prefix, func(ctx *context.Context) {
		ctx.Values().Set("code", 1)
		ctx.Next()
	})

	//// mysql config
	//mc := MysqlConfig{
	//	Host:     "127.0.0.1",
	//	Port:     3306,
	//	Username: "test",
	//	Password: "testPassword",
	//	DbName:   "test",
	//	PoolSize: 100,
	//	ShowSql:  true,
	//}
	// mysql instance
	mdb, _ := xorm.NewEngine("sqlite3", "./test.db")
	_ = mdb.Sync2(new(testModel))
	mdb.ShowSQL(true)
	//// redis config
	//rc := RedisConfig{
	//	Host:     "127.0.0.1",
	//	Port:     6379,
	//	Password: "123456789",
	//	Db:       6,
	//	PoolSize: 100,
	//}
	// redis instance
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456789",
		DB:       5,
	})

	checkMc := &Config{
		Party: p,
		MysqlInstance: MysqlInstance{
			Mdb: mdb,
		},
		RedisInstance: RedisInstance{
			Rdb: rdb,
		},
		Models: []*SingleModel{
			{
				Model:     new(testModel),
				CacheTime: 1 * time.Minute,
				//PrivateContextKey: "code",
				//PrivateColName:    "code",
				//GetAllExtraFilters: map[string]string{
				//	"code": "1",
				//},
				//GetSingleExtraFilters: map[string]string{
				//	"code": "2",
				//},

			},
		},
	}
	New(checkMc)
	testModel := mdb.TableName(checkMc.Models[0].Model)
	fp := prefix + "/" + testModel
	e := httptest.New(t, app)
	testCrud(t, e, fp)
	// because use delay delete to default 500ms
	time.Sleep(600 * time.Millisecond)
	testCache(t, e, fp)
}

// test crud
func testCrud(t *testing.T, e *httpexpect.Expect, fp string) {
	println("run crud test")
	// get all
	getAll := e.GET(fp).Expect().Status(httptest.StatusOK)
	getAll.JSON().Object().ContainsKey("page")

	// add new data
	bodyMap := map[string]interface{}{"name": "test", "age": 68, "desc": "desc"}
	addData := e.POST(fp).WithForm(bodyMap).Expect().Status(httptest.StatusOK)
	addData.JSON().Object().Value("name").Equal("test")
	id := addData.JSON().Object().Value("id").Raw()
	println("get data list")
	fs := fp + "/" + fmt.Sprintf("%v", id)

	// get single
	getSingle := e.GET(fs).Expect().Status(httptest.StatusOK)
	getSingle.JSON().Object().ContainsKey("name")
	println("get single data")

	// put data
	editMap := map[string]interface{}{"name": "edit"}
	edit := e.PUT(fs).WithForm(editMap).Expect().Status(httptest.StatusOK)
	edit.JSON().Object().Value("name").Equal("edit")
	println("put data")

	// delete data
	deleteData := e.DELETE(fs).Expect().Status(httptest.StatusOK)
	deleteData.JSON().Object().ContainsKey("id")
	println("delete data")
}

// test cache
func testCache(t *testing.T, e *httpexpect.Expect, fp string) {
	println("run cache test")

	fs := fp + "/4"
	// get all save to redis
	e.GET(fp).Expect()

	// get redis cache
	cacheAll := e.GET(fp).Expect().Status(httptest.StatusOK)
	cacheAll.JSON().Object().Value("status").Equal("cache")
	println("cache all data list")

	e.GET(fs).Expect()
	cacheSingle := e.GET(fs).Expect().Status(httptest.StatusOK)
	cacheSingle.JSON().Object().Value("status").Equal("cache")
	println("cache single data")
}
