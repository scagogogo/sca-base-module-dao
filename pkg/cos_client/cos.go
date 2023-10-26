package cos_client

import (
	"bytes"
	"context"
	"fmt"
	sca_base_module_config "github.com/scagogogo/sca-base-module-config"
	"github.com/scagogogo/sca-base-module-logger/logger"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// ------------------------------------------------ ACL ----------------------------------------------------------------

// ObjectACL 到底是我傻逼还是腾讯云傻逼？它的SDK里都是直接用的字符串？还是搞个枚举吧先...等发现有枚举了再替换
type ObjectACL string

const (

	// ObjectACLPublicRead 公共读取
	ObjectACLPublicRead = "public-read"

	// ObjectACLPublicWrite 公共读写
	ObjectACLPublicWrite = "public-read-write"

	// ObjectACLPrivate 私有
	ObjectACLPrivate = "private"
)

// -------------------------------------------- 基础方法 ---------------------------------------------------------------

// New 创建一个COS客户端，暂不支持在创建时指定bucket和区域，bucket和区域在配置文件中指定，以后有需要再提供指定参数的支持
func New() (*cos.Client, error) {

	if sca_base_module_config.Config == nil {
		return nil, fmt.Errorf("config module not init")
	}

	//将<bucket>和<region>修改为真实的信息
	//bucket的命名规则为{name}-{appid} ，此处填写的存储桶名称必须为此格式
	// "https://<bucket>.cos.<region>.myqcloud.com"
	u, _ := url.Parse(sca_base_module_config.Config.COS.Endpoint)
	b := &cos.BaseURL{BucketURL: u}
	return cos.NewClient(b, &http.Client{
		//设置超时时间
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			//如实填写账号和密钥，也可以设置为环境变量
			SecretID:  sca_base_module_config.Config.COS.SecretId,
			SecretKey: sca_base_module_config.Config.COS.SecretKey,
			// Debug 模式，把对应 请求头部、请求内容、响应头部、响应内容 输出到标准输出
			//Transport: &debug.DebugRequestTransport{
			//	RequestHeader:  true,
			//	RequestBody:    true,
			//	ResponseHeader: true,
			//	ResponseBody:   false,
			//},
		},
	}), nil
}

// -------------------------------------------- 开箱即用方法 ------------------------------------------------------------

// PutPublicObject 把字节对象上传到oss，权限是公开读
// bucketName: 要上传到的bucket的名字
// objectName: 对象的路径
// objectBytes: 对象的数据
func PutPublicObject(ctx context.Context,  c *cos.Client, objectName string, objectBytes []byte) error {
	return PutObject(ctx, c, objectName, objectBytes, ObjectACLPublicRead)
}

// PutPrivateObject 把字节对象上传到oss，权限是私有读写
// bucketName: 要上传到的bucket的名字
// objectName: 对象的路径
// objectBytes: 对象的数据
func PutPrivateObject(ctx context.Context,  c *cos.Client, objectName string, objectBytes []byte) error {
	return PutObject(ctx, c, objectName, objectBytes, ObjectACLPrivate)
}

// PutObject 把字节对象上传到oss
// objectName: 对象的路径
// objectBytes: 对象的数据
// ossAcl: 对象要设定的权限
func PutObject(ctx context.Context, c *cos.Client, objectName string, objectBytes []byte, objectAcl ObjectACL) error {

	opt := &cos.ObjectPutOptions{
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: string(objectAcl),
		},
	}

	r, err := c.Object.Put(ctx, objectName, bytes.NewReader(objectBytes), opt)
	if err != nil {
		return err
	}
	// 只是简单的判断一下http的响应码
	if r.StatusCode != 200 {
		return fmt.Errorf("put object %s response statuc code error: %d", objectName, r.StatusCode)
	}

	return nil
}

// GetObjectBytes 获取指定路径的OSS对象的数据，以字节数组的形式返回
// bucketName: 对象所在的bucket的名字
// objectName: 对象的路径（名字）
func GetObjectBytes(ctx context.Context, c *cos.Client, objectName string) ([]byte, error) {

	// 取出对象
	resp, err := c.Object.Get(ctx, objectName, nil)
	if err != nil {
		return nil, err
	}

	objectBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		logger.Errorf("objectName = %s, close GetObject ReadCloser error：%s", objectName, err.Error())
	}

	return objectBytes, nil
}


