package main

import (
	"bytes"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"gopkg.in/mgo.v2"
)

var (
	session *mgo.Session
)

// DBOut 数据库输出
type DBOut struct {
	BsonString string
	Type       string
	Name       string
}

func connectDB(dbURL string) {
	var err error
	session, err = mgo.Dial(dbURL)
	if err != nil {
		panic("Mongo Dial Err: " + err.Error() + "\r\nAtUrl: " + dbURL)
	}
	session.SetMode(mgo.Monotonic, true)
}

func main() {
	var dbhost string
	var dbport int
	var dburl string
	flag.StringVar(&dbhost, "s", "127.0.0.1", "mongo所在host")
	flag.IntVar(&dbport, "p", 27017, "mongo端口")
	flag.StringVar(&dburl, "u", "", "host:port")
	flag.Parse()

	if len(dburl) > 0 {
		connectDB(dburl)
	} else {
		connectDB(dbhost + ":" + strconv.Itoa(dbport))
	}

	mongo := session.Clone()
	defer mongo.Close()

	dbs, err := mongo.DatabaseNames()
	if err != nil {
		panic("Find DBS Err:" + err.Error())
	}
	for _, db := range dbs {
		doDataBase(db)
	}
}

func doDataBase(dbName string) {
	mongo := session.Clone()
	defer mongo.Close()

	cs, err := mongo.DB(dbName).CollectionNames()
	if err != nil {
		return
	}
	for _, c := range cs {
		doCollection(dbName, c)
	}
}

func doCollection(dbName, collection string) {
	mongo := session.Clone()
	defer mongo.Close()

	resultMap := map[string]DBOut{}
	var query map[string]interface{}

	iter := mongo.DB(dbName).C(collection).Find(nil).Iter()
	for iter.Next(&query) {
		makeQueryKeyAndType(resultMap, query)
	}

	structString, _ := parse(resultMap, collection)
	fmt.Println(structString)
}

func makeQueryKeyAndType(result map[string]DBOut, query map[string]interface{}) {
	ifFirstIn := len(result) == 0
	hadExist := map[string]bool{}
	for key, bsonValue := range query {
		hadExist[key] = true
		typeString := strings.Replace(fmt.Sprintf("%T", bsonValue), " ", "", -1)

		if _, ok := result[key]; ok {
			//存在
		} else if !ifFirstIn {
			//不存在，且已经检测过result
			result[key] = DBOut{
				BsonString: fmt.Sprintf("`bson:\"%v,omitempty\"`", key), //type
				Type:       typeString,
				Name:       key,
			}
		} else {
			//首次进来
			result[key] = DBOut{
				BsonString: fmt.Sprintf("`bson:\"%v\"`", key), //type
				Type:       typeString,
				Name:       key,
			}
		}
	}

	// 查询是否有减少字段
	// 规范Name
	for key, item := range result {
		name := strings.ToUpper(key[0:1])
		if len(key) > 1 {
			name += key[1:]
		}
		if name == "_id" {
			name = "ID"
		}
		if ok := hadExist[key]; !ok {
			result[key] = DBOut{
				BsonString: fmt.Sprintf("`bson:\"%v,omitempty\"`", key), //type
				Type:       item.Type,
				Name:       name,
			}
		} else {
			result[key] = DBOut{
				BsonString: item.BsonString, //type
				Type:       item.Type,
				Name:       name,
			}
		}
	}
}

func parse(query map[string]DBOut, structName string) (string, error) {
	const templ = `{
{{range $key, $item := .}}	{{$item.Name}}	{{$item.Type}}	{{$item.BsonString}}
{{end}}}
`
	var doc bytes.Buffer
	t := template.Must(template.New("struct").Parse(templ))
	if err := t.Execute(&doc, query); err != nil {
		return "", err
	}
	return fmt.Sprintf("type %s struct%s", structName, doc.String()), nil
}
