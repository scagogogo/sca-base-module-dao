package oss_client

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	sca_base_module_config "github.com/scagogogo/sca-base-module-config"
	"github.com/scagogogo/sca-base-module-logger/logger"
	"io/ioutil"
	"sync"
)

// -------------------------------------------- 基础方法 ---------------------------------------------------------------

// NewOssClient 创建一个OSS客户端，bucket的名字可以自己指定
func NewOssClient(bucketName string) (*oss.Bucket, error) {
	if sca_base_module_config.Config == nil {
		return nil, fmt.Errorf("config module not init")
	}
	ossConfig := sca_base_module_config.Config.OSS
	client, err := oss.New(ossConfig.Endpoint, ossConfig.AccessKeyId, ossConfig.AccessKeySecret)
	if err != nil {
		logger.Errorf("bucketName = %s, create oss client error: %s", bucketName, err.Error())
		return nil, err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		logger.Errorf("bucketName = %s, create oss client error: %s", bucketName, err.Error())
		return nil, err
	}

	return bucket, nil
}

// -------------------------------------------- 开箱即用方法 ------------------------------------------------------------

// 存储bucket name对应的bucket
var bucketMap map[string]*oss.Bucket
var bucketMapLock = sync.Mutex{}

func takeBucket(bucketName string) (*oss.Bucket, error) {

	// 这个方法操作bucket是互斥的
	bucketMapLock.Lock()
	defer bucketMapLock.Unlock()

	// 如果bucketMap不存在的话就先初始化它
	if bucketMap == nil {
		bucketMap = make(map[string]*oss.Bucket, 0)
	}

	// 取出这个bucket对应的客户端
	bucket := bucketMap[bucketName]
	if bucket != nil {
		return bucket, nil
	}

	// 如果之前没初始化过的话，先初始化
	var err error
	bucket, err = NewOssClient(bucketName)
	if err != nil {
		return nil, err
	}
	bucketMap[bucketName] = bucket

	return bucket, nil
}

// PutPublicObjectToOss 把字节对象上传到oss，权限是公开读
// bucketName: 要上传到的bucket的名字
// objectName: 对象的路径
// objectBytes: 对象的数据
func PutPublicObjectToOss(bucketName, objectName string, objectBytes []byte) error {
	return PutObjectToOss(bucketName, objectName, objectBytes, oss.ACLPublicRead)
}

// PutPrivateObjectToOss 把字节对象上传到oss，权限是私有读写
// bucketName: 要上传到的bucket的名字
// objectName: 对象的路径
// objectBytes: 对象的数据
func PutPrivateObjectToOss(bucketName, objectName string, objectBytes []byte) error {
	return PutObjectToOss(bucketName, objectName, objectBytes, oss.ACLPrivate)
}

// PutObjectToOss 把字节对象上传到oss
// bucketName: 要上传到的bucket的名字
// objectName: 对象的路径
// objectBytes: 对象的数据
// ossAcl: 对象要设定的权限
func PutObjectToOss(bucketName, objectName string, objectBytes []byte, ossAcl oss.ACLType) error {

	// 获取到对应的Bucket
	bucket, err := takeBucket(bucketName)
	if err != nil {
		return err
	}

	// 上传到OSS上
	// 指定存储类型为标准存储，缺省也为标准存储。
	storageType := oss.ObjectStorageClass(oss.StorageStandard)
	// 指定访问权限为公共读，缺省为继承bucket的权限。
	objectAcl := oss.ObjectACL(ossAcl)

	// 尝试上传对象到OSS上
	err = bucket.PutObject(objectName, bytes.NewReader(objectBytes), storageType, objectAcl)
	if err != nil {
		return err
	}
	return nil
}

// GetObjectBytes 获取指定路径的OSS对象的数据，以字节数组的形式返回
// bucketName: 对象所在的bucket的名字
// objectName: 对象的路径（名字）
func GetObjectBytes(bucketName, objectName string) ([]byte, error) {

	// 获取到对应的Bucket
	bucket, err := takeBucket(bucketName)
	if err != nil {
		return nil, err
	}

	// 取出对象
	objectReader, err := bucket.GetObject(objectName)
	if err != nil {
		return nil, err
	}

	objectBytes, err := ioutil.ReadAll(objectReader)
	if err != nil {
		return nil, err
	}
	err = objectReader.Close()
	if err != nil {
		logger.Errorf("bucketName = %s, objectName = %s, close GetObject ReadCloser error：%s", bucketName, objectName, err.Error())
	}

	return objectBytes, nil
}

// --------------------------------------------------------------------------------------------------------------------
