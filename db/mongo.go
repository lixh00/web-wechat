package db

import (
	"context"
	"gitee.ltd/lxh/logger/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"web-wechat/core"
)

type mongoDBClient struct {
	client *mongo.Client
}

var MongoClient mongoDBClient

// InitMongoConnHandle 初始化MongoDB连接
func InitMongoConnHandle() {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // 在调用WithTimeout之后defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(core.SystemConfig.MongoDbConfig.GetClientUri()))
	if err != nil {
		log.Panicf("MongoDB初始化连接失败: %v", err.Error())
		//os.Exit(1)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Panicf("MongoDB初始化连接失败: %v", err.Error())
	}

	log.Info("MongoDB连接初始化成功")
	//mongoClient = client
	MongoClient = mongoDBClient{client: client}
}

// Save 保存数据到Mongo
func (m *mongoDBClient) Save(data interface{}, tableName string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // 在调用WithTimeout之后defer cancel()

	collection := m.client.Database(core.SystemConfig.MongoDbConfig.DbName).Collection(tableName)
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		log.Errorf("保存数据到MongoDB失败: %v", err.Error())
		return false
	}
	log.Debugf("MongoDB保存结果: %v", res)
	return true
}
