package main

import (
	"github.com/23233/ab"
	"github.com/23233/ab/_examples/model"
	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var Engine *xorm.Engine
var Rdb *redis.Client

func init() {
	// database 连接器
	var err error

	Engine, err = xorm.NewEngine("sqlite3", "./simple.db")

	if err != nil {
		println(err.Error())
		return
	}
	//Engine.SetLogger()
	Engine.ShowSQL(true)
	//Engine.ShowExecTime(true)
	err = Engine.Ping()
	if err != nil {
		panic(err)
	}

	// redis instance
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "OyYxP4HCX9dNuMKkORHWH9vhHoJJNoti",
		DB:       5,
	})

}

func NewApp() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel("debug")

	modelList := []interface{}{
		new(model.TestModelA),
		new(model.TestModelB),
		new(model.ComplexModelC),
		new(model.ComplexModelD),
		new(model.TestStructComplexModel),
	}

	_ = Engine.Sync2(modelList...)

	app.Get("/", func(ctx iris.Context) {
		_, _ = ctx.JSON(iris.Map{})
	})

	v1 := app.Party("/api/v1")
	abConfig := ab.Config{
		Party: v1,
		MysqlInstance: ab.MysqlInstance{
			//MysqlConfig: mc,
			Mdb: Engine,
		},
		RedisInstance: ab.RedisInstance{
			//RedisConfig: rc,
			Rdb: Rdb,
		},
		StructList: []ab.SingleModel{
			{
				Model:             new(model.TestModelA),
				GetAllResponse:    new(model.TestModelResp),
				GetSingleResponse: new(model.TestModelResp),
				PostResponse:      new(model.TestModelResp),
				PutResponse:       new(model.TestModelResp),
			},
			{
				Model: new(model.TestModelB),
			},
			{
				Model: new(model.ComplexModelC),
			},
			{
				Model: new(model.ComplexModelD),
			},
			{
				Model: new(model.TestStructComplexModel),
			},
		},
	}
	ab.New(&abConfig)
	return app
}

func main() {
	app := NewApp()
	_ = app.Listen(":8080")
}
