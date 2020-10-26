package main

import (
	"github.com/23233/ab"
	"github.com/23233/ab/_examples/model"
	"github.com/kataras/iris/v12"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var Engine *xorm.Engine

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
}

func main() {
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

	api := ab.New(ab.Config{
		App:        app,
		StructList: modelList,
		Engine:     Engine,
		Prefix:     "/api/v1/",
	})
	api.Run()

	_ = app.Listen(":8080")

}