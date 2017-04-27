package database

import (
	"time"

	mgo "gopkg.in/mgo.v2"
)

// mongoDB MongoDB实例
type mongoDB struct {
	sess *mgo.Session
}

// InitDB 初始化
func (d *mongoDB) InitDB(uri string) error {
	var err error
	sess, err := mgo.Dial(uri)
	if err != nil {
		return err
	}

	sess.SetMode(mgo.Monotonic, true)
	sess.SetSocketTimeout(2 * time.Minute)
	d.sess = sess
	return nil
}

// GetCollection GetCollection
func (d *mongoDB) GetCollection(config Configer) (*mgo.Collection, CloseSessionFunc) {
	if d == nil {
		return nil, nil
	}
	tmp := d.sess.Clone()
	tmp.SetSocketTimeout(time.Duration(config.GetSocketTimeoutSecond()) * time.Second)
	return tmp.DB(config.GetDBName()).C(config.GetCollectionName()), func() { tmp.Close() }
}
