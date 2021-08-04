package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"web-wechat/core"
	"web-wechat/logger"
)

var mongoClient *mongo.Client

// InitMongoConnHandle 初始化MongoDB连接
func InitMongoConnHandle() {
	// 读取配置
	core.InitMongoConfig()
	client, err := mongo.NewClient(options.Client().ApplyURI(core.MongoDbConfig.GetClientUri()))
	if err != nil {
		logger.Log.Panicf("MongoDB初始化连接失败: %v", err.Error())
		//os.Exit(1)
	} else {
		logger.Log.Info("MongoDB连接初始化成功")
		mongoClient = client
	}
}

// SaveToMongo 保存数据到Mongo
func SaveToMongo(data interface{}, tableName string) bool {
	ctx := context.Background()
	collection := mongoClient.Database(core.MongoDbConfig.DbName).Collection(tableName)
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		logger.Log.Errorf("保存数据到MongoDB失败: %v", err.Error())
		return false
	}
	logger.Log.Debugf("MongoDB保存结果: %v", res)
	return true
}
