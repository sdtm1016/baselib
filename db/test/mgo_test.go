package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"baselib/logger"
)

type Address struct {
	Address string
}

type Location struct {
	Longitude float64
	Latitude  float64
}

type Person struct {
	Id       bson.ObjectId `bson:"_id"`
	Name     string
	Age_Int  int
	Address  []Address
	Location Location
}

func main() {
	logger.Info("start to connect to mongoDB...")
	//连接
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		logger.Error(err)
		return
	}
	//设置模式
	session.SetMode(mgo.Monotonic, true)
	//获取文档集
	collection := session.DB("test").C("person")
	// 创建索引
	index := mgo.Index{
		Key:        []string{"name"}, // 索引字段， 默认升序,若需降序在字段前加-
		Unique:     true,             // 唯一索引 同mysql唯一索引
		DropDups:   true,             // 索引重复替换旧文档,Unique为true时失效
		Background: true,             // 后台创建索引
	}
	if err := collection.EnsureIndex(index); err != nil {
		logger.Error(err)
		return
	}
	// 创建一个范围索引
	if err := collection.EnsureIndexKey("$2dsphere:location"); err != nil {
		logger.Error(err)
		return
	}
	//添加记录
	person := Person{
		Id:      bson.NewObjectId(),
		Name:    "逍遥子",
		Age_Int: 24,
		Address: []Address{
			Address{
				Address: "逍遥谷",
			},
		},
		Location: Location{
			Longitude: 1,
			Latitude:  1,
		},
	}
	if err := collection.Insert(person); err != nil {
		logger.Error(err)
		return
	}

	//查找记录
	newPerson := &Person{}
	if err := collection.Find(bson.M{"age_int": 24}).One(newPerson); err != nil {
		logger.Error(err)
		return
	}

	//修改记录
	if err := collection.Update(bson.M{"age_int": 24}, bson.M{"$set": bson.M{"age_int": 26}}); err != nil {
		logger.Error(err)
		return
	}

	//删除记录
	//if err := collection.Remove(bson.M{"age_int": 26}); err != nil {
	//  logger.Error(err)
	//  return
	//}

	//位置搜索
	selector := bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{1, 1},
				},
				"$maxDistance": 1,
				//"$minDistance": 0,
			},
		},
	}

	if err := collection.Find(selector).One(newPerson); err != nil {
		logger.Error(err)
		return
	}

	//close
	session.Close()
}