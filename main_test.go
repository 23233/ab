package ab

import (
	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12"
	_ "github.com/mattn/go-sqlite3"
	"testing"
	"xorm.io/xorm"
)

type testModel struct {
	Name string `xorm:"varchar(10)" json:"name"`
	Age  uint64 `json:"age"`
	Desc string `xorm:"varchar(20)" json:"desc"`
}

func TestNew(t *testing.T) {
	app := iris.New()

	p := app.Party("/")

	// mysql config
	mc := MysqlConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "test",
		Password: "testPassword",
		DbName:   "test",
		PoolSize: 100,
		ShowSql:  true,
	}
	// mysql instance
	mdb, _ := xorm.NewEngine("sqlite3", "./test.db")
	// redis config
	rc := RedisConfig{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "123456789",
		Db:       6,
		PoolSize: 100,
	}
	// redis instance
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456789",
		DB:       5,
	})

	// test msql config valid
	checkMc := &Config{
		Party: p,
		MysqlInstance: MysqlInstance{
			MysqlConfig: mc,
			Mdb:         mdb,
		},
		RedisInstance: RedisInstance{
			RedisConfig: rc,
			Rdb:         rdb,
		},
		StructList: []SingleModel{
			{
				Model: new(testModel),
			},
		},
	}
	New(checkMc)
}
