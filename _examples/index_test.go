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
	// 测试新增
	bodyMap := map[string]interface{}{"name": "test"}
	addData := e.POST(prefix + "/" + testModel).WithForm(bodyMap).Expect().Status(httptest.StatusOK)
	addData.JSON().Object().Empty()

}
