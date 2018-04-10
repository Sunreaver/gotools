package database

import (
	"reflect"
	"testing"

	. "github.com/bouk/monkey"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	mgo "gopkg.in/mgo.v2"

	mock_database "github.com/sunreaver/gotools/database/mock"
)

func TestInitMongoDB(t *testing.T) {
	Convey("test init mongo db", t, func() {
		Convey("for error", func() {
			ctl := gomock.NewController(t)
			mockConfiger := mock_database.NewMockConfiger(ctl)
			mockConfiger.EXPECT().GetDBName().Return("mongo")
			mockConfiger.EXPECT().GetCollectionName().Return("collection")
			mockConfiger.EXPECT().GetSocketTimeoutSecond().Return(10)

			defer UnpatchAll()
			var mgoDB *mongoDB
			PatchInstanceMethod(reflect.TypeOf(mgoDB), "InitDB", func(_ *mongoDB, _ string) error {
				return nil
			})
			PatchInstanceMethod(reflect.TypeOf(mgoDB), "GetCollection", func(_ *mongoDB, cfg Configer) (*mgo.Collection, CloseSessionFunc) {
				So(cfg.GetDBName(), ShouldEqual, "mongo")
				So(cfg.GetCollectionName(), ShouldEqual, "collection")
				So(cfg.GetSocketTimeoutSecond(), ShouldEqual, 10)
				return nil, func() {}
			})

			e := InitMongoDB("test1")
			So(e, ShouldBeNil)
			Mongo().GetCollection(mockConfiger)
		})
	})
}
