package main

import (
	"github.com/kataras/iris/v12/httptest"
	"testing"
)

func TestNewApp(t *testing.T) {
	app := NewApp()
	e := httptest.New(t, app)
	prefix := "/api/v1"
	testModel := "test_model_a"
	//// 测试新增 不支持withjson 因为iris只有readJson 没有单个json fields的获取
	//bodyMap := map[string]interface{}{"name": "test"}
	//addData := e.POST(prefix + "/" + testModel).WithForm(bodyMap).Expect().Status(httptest.StatusOK)
	//addData.JSON().Object().ContainsKey("name").Equal("test")
	//// 测试获取
	//getAll := e.GET(prefix + "/" + testModel).Expect().Status(httptest.StatusOK)
	//getAll.JSON().Object().ContainsKey("page")

	// 测试获取单个
	getSingle := e.GET(prefix + "/" + testModel + "/6").Expect().Status(httptest.StatusOK)
	getSingle.JSON().Object().ContainsKey("name")

	// 测试修改
	editMap := map[string]interface{}{"name": "修改后"}
	editData := e.PUT(prefix + "/" + testModel + "/6").WithForm(editMap).Expect().Status(httptest.StatusOK)
	editData.JSON().Object().Equal(editMap)

	// 测试删除
	deleteData := e.DELETE(prefix + "/" + testModel + "/5").Expect().Status(httptest.StatusOK)
	deleteData.JSON().Object().ContainsKey("id")

}
