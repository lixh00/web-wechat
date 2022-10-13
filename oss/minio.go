package oss

import (
	"context"
	"gitee.ltd/lxh/logger/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"web-wechat/core"
)

var minioClient *minio.Client

// InitOssConnHandle 初始化OSS连接
func InitOssConnHandle() {
	ctx := context.Background()
	// 初使化 minio client对象。
	client, err := minio.New(core.SystemConfig.OssConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(core.SystemConfig.OssConfig.AccessKeyID, core.SystemConfig.OssConfig.SecretAccessKey, ""),
		Secure: core.SystemConfig.OssConfig.UseSsl,
	})
	if err != nil {
		log.Panicf("OSS初始化失败: %v", err.Error())
	}
	log.Info("OSS连接成功，开始判断桶是否存在")
	// 判断捅是否存在，不存在就创建
	exists, err := client.BucketExists(ctx, core.SystemConfig.OssConfig.BucketName)
	if err != nil {
		log.Panicf("判断桶失败: %v", err)
	}
	if !exists {
		log.Info("桶不存在，开始创建")
		// 创建桶
		err = client.MakeBucket(ctx, core.SystemConfig.OssConfig.BucketName, minio.MakeBucketOptions{Region: "us-east-1"})
		if err != nil {
			log.Panicf("OSS桶创建失败: %v", err.Error())
		}
		log.Info("桶创建成功")
	} else {
		log.Info("桶已存在")
	}
	minioClient = client
	log.Info("OSS初始化成功")
}

// SaveToOss 保存文件到OSS
func SaveToOss(b io.Reader, contentType, fileName string) bool {
	log.Debugf("开始上传文件: %v", fileName)
	ctx := context.Background()
	_, err := minioClient.PutObject(ctx, core.SystemConfig.OssConfig.BucketName, fileName, b, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Errorf("文件上传错误: %v", err)
		return false
	}
	log.Debugf("文件上传完毕: %v", fileName)
	return true
}
