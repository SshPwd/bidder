package db

import (
	"fmt"
	"golang_bidder/config"
	"golang_bidder/logs"
	"gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson"
	"os"
)

var (
	globalMgoSession *mgo.Session
)

func InitMongodb() {

	conf := config.Get()

	url := ""

	if conf.Mongodb.User != "" {

		url = fmt.Sprintf("mongodb://%s:%s@%s/%s",
			conf.Mongodb.User,
			conf.Mongodb.Password,
			conf.Mongodb.Host,
			conf.Mongodb.DbName)

	} else {

		url = fmt.Sprintf("mongodb://%s/%s",
			conf.Mongodb.Host,
			conf.Mongodb.DbName)
	}

	session, err := mgo.Dial(url)
	if err != nil {

		logs.Critical(fmt.Sprintf("Mongodb (%s/%s): %s",
			conf.Mongodb.Host, conf.Mongodb.DbName, err.Error()))

		os.Exit(1)
		return
	}

	// session.SetMode(mgo.Monotonic, true)

	err = session.Ping()
	if err != nil {

		logs.Critical(fmt.Sprintf("Mongodb (%s/%s): %s",
			conf.Mongodb.Host, conf.Mongodb.DbName, err.Error()))

		os.Exit(1)
		return
	}

	// session.DB(conf.Mongodb.DbName).C(conf.Mongodb.CollectionDspTitle).EnsureIndex(mgo.Index{
	// 	Key:        []string{"dsp_id", "date", "title_hash"},
	// 	Unique:     true,
	// 	DropDups:   true,
	// 	Background: false,
	// 	Sparse:     false,
	// })

	// session.DB(conf.Mongodb.DbName).C(conf.Mongodb.CollectionDsp).EnsureIndex(mgo.Index{
	// 	Key:        []string{"dsp_id", "date"},
	// 	Unique:     true,
	// 	DropDups:   true,
	// 	Background: false,
	// 	Sparse:     false,
	// })

	globalMgoSession = session

	fmt.Println("InitMongodb:", globalMgoSession)
}
