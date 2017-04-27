package database

import mgo "gopkg.in/mgo.v2"

// CloseSessionFunc 关闭sess Func
type CloseSessionFunc func()

// Configer Configer
type Configer interface {
	GetDBName() string
	GetCollectionName() string
	GetSocketTimeoutSecond() int
}

// MgoExamples mongo实例
type MgoExamples interface {
	GetCollection(config Configer) (*mgo.Collection, CloseSessionFunc)
}

var (
	// Mongo 服务中的Mongo实例
	mongo MgoExamples
)

// InitMongoDB InitMongoDB
func InitMongoDB(uri string) error {
	var db = mongoDB{}
	e := db.InitDB(uri)
	if e != nil {
		return e
	}
	mongo = &db
	return nil
}

// Mongo 获取mongo实例。唯一对外获取实例接口。
func Mongo() MgoExamples {
	return mongo
}
