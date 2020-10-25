package model

import (
	"time"
)

type ModelBase struct {
	Id         uint64    `xorm:"autoincr pk unique" json:"id"`
	CreateTime time.Time `xorm:"created index" json:"create_time"`
	UpdateTime time.Time `xorm:"updated" json:"update_time"`
	DeletedAt  time.Time `xorm:"deleted" json:"deleted_at"`
	Version    uint16    `xorm:"version" json:"version"`   //版本号
	Status     uint8     `xorm:"default(0)" json:"status"` // 当前状态 0 正常 其他都不正常
}

type TestModelA struct {
	Id   uint64 `xorm:"autoincr pk unique" json:"id"`
	Name string `xorm:"varchar(20)" comment:"名称"`
}

func (u *TestModelA) SpAlias() string {
	return "赵日天"
}

type CustomReqValid struct {
	Id   uint64 `json:"id" sp:"lineTo(Id)"`
	Uid  string `json:"uid" comment:"标识符" validate:"required" sp:"lineTo(Uid)"`
	Name string `json:"name" comment:"名称" validate:"required"`
}

type CustomReqBValid struct {
	Desc string `json:"desc" validate:"required"`
}

type TestModelB struct {
	Id   uint64 `xorm:"autoincr pk unique" json:"id"`
	Desc string `xorm:"varchar(60)"`
}

type ComplexModelC struct {
	Id      uint64 `xorm:"autoincr pk unique" json:"id"`
	Name    string `xorm:"varchar(20)" json:"name"`
	NowTime time.Time
	Count   uint
}

type ComplexModelD struct {
	Id               uint64        `xorm:"autoincr pk unique" json:"id"`
	Name             string        `xorm:"varchar(20)" json:"name"`
	TestString       string        `xorm:"varchar(20)" json:"test_string"`
	TestInt          int           `xorm:"" json:"test_int"`
	TestInt8         int8          `xorm:"" json:"test_int_8"`
	TestInt16        int16         `xorm:"" json:"test_int_16"`
	TestInt32        int32         `xorm:"" json:"test_int_32"`
	TestInt64        int64         `xorm:"" json:"test_int_64"`
	TestUint         uint          `xorm:"" json:"test_uint"`
	TestUint8        uint8         `xorm:"" json:"test_uint_8"`
	TestUint16       uint16        `xorm:"" json:"test_uint_16"`
	TestUint32       uint32        `xorm:"" json:"test_uint_32"`
	TestUint64       uint64        `xorm:"" json:"test_uint_64"`
	TestFloat32      float32       `json:"test_float_32"`
	TestFloat64      float64       `json:"test_float_64"`
	TestTimeDuration time.Duration `json:"test_time_duration"`
	TestTimeTime     time.Time     `json:"test_time_time"`
	TestBool         bool          `json:"test_bool"`
	TestFk           string        `json:"test_fk" sp:"fk('ComplexModelC')"`            // one to one
	TestManyFk       string        `json:"test_many_fk" sp:"fk:ComplexModelC multiple"` // one to many
}

func (c *ComplexModelD) ApiSearchFields() []string {
	return []string{"Name", "TestString"}
}

type TestUserModel struct {
	Id       uint64 `xorm:"autoincr pk unique" json:"id"`
	UserName string `xorm:"varchar(60) notnull" json:"user_name"`
	Password string `xorm:"varchar(100) notnull" json:"password"`
	Salt     string `xorm:"varchar(40) notnull" json:"salt"`
	Niubi    string `xorm:"varchar(30)"`
}

type TestStructComplexModel struct {
	ModelBase `xorm:"extends"`
	Names     string `xorm:"notnull" json:"names"`
	Img       string `xorm:"varchar(255)" sp:"tag(img)"`
}
