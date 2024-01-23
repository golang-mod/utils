package alioss

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Config struct {
	AccessKeyId     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	Bucket          string `yaml:"bucket"`
	Endpoint        string `yaml:"endpoint"`
	EndpointUrl     string `yaml:"endpoint_url"` // 格式为 bucketname.endpoint
	CallbackUrl     string `yaml:"callback_url"` // 回调地址
	UploadDir       string `yaml:"upload_dir"`   // 用户上传文件时指定的前缀
	ExpireTime      int64  `yaml:"expire_time"`  // 超时时间 default:30
}

type OSS struct {
	Config  Config
	TmpPath string // 临时文件目录
	client  *oss.Client
	bucket  *oss.Bucket
	done    uint32
	m       sync.Mutex
}

func New(c Config) *OSS {
	var o = new(OSS)
	o.Config = c
	o.TmpPath = "tmp/"
	return o
}

func Version() {
	fmt.Println("OSS Go SDK Version: ", oss.Version)
}

func (o *OSS) connection() error {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done != 0 {
		return nil
	}
	o.done = 1
	fmt.Println("连接OSS")
	// 创建OSSClient实例。
	client, err := oss.New(o.Config.Endpoint, o.Config.AccessKeyId, o.Config.AccessKeySecret)
	if err != nil {
		return err
	}
	o.client = client
	// 获取存储空间。
	bucket, err := client.Bucket(o.Config.Bucket)
	if err != nil {
		return err
	}
	o.bucket = bucket
	// 指定存储类型为标准存储，缺省也为标准存储。
	//storageType := oss.ObjectStorageClass(oss.StorageStandard)
	// 指定存储类型为归档存储。
	// storageType := oss.ObjectStorageClass(oss.StorageArchive)
	// 指定访问权限为公共读，缺省为继承bucket的权限。
	//objectAcl := oss.ObjectACL(oss.ACLPublicRead)
	return nil
}

// Signature 阿里云OSS签名上传
func (o *OSS) Signature() PolicyToken {
	return o.getPolicyToken()
}

// SignUrl 私有文件签名访问
func (o *OSS) SignUrl(objectKey string) (url string, err error) {
	err = o.connection()
	if err != nil {
		return "", err
	}
	signURL, err := o.bucket.SignURL(objectKey, "GET", o.Config.ExpireTime*100)
	if err != nil {
		return "", err
	}
	return signURL, nil
}

// ModifiedTime 获取文件变更时间
func (o *OSS) ModifiedTime(objectKey string) (datetime string, err error) {
	err = o.connection()
	if err != nil {
		return
	}
	meta, err := o.bucket.GetObjectMeta(objectKey)
	if err != nil {
		return
	}
	date := meta.Get("Last-Modified")
	onlineAt, err := time.ParseInLocation(time.RFC1123, date, time.UTC)
	return onlineAt.Local().Format(time.DateTime), nil
}

// CopyTmpFile 拷贝tmp临时文件
// tmpUrls 临时文件列表 无"tmp/"前缀不拷贝
// copyUrls 拷贝后文件列表 无"tmp/"前缀不拷贝
func (o *OSS) CopyTmpFile(tmpUrls []string) (copyUrls []string, err error) {
	err = o.connection()
	if err != nil {
		return copyUrls, err
	}
	if err != nil {
		return copyUrls, err
	}
	for _, u := range tmpUrls {
		uParse, err := url.Parse(u)
		if err != nil {
			return copyUrls, err
		}
		srcObjectKey := uParse.Path
		if strings.HasPrefix(srcObjectKey, "/") {
			srcObjectKey = strings.Replace(srcObjectKey, "/", "", 1)
			//srcObjectKey = srcObjectKey[1:len(srcObjectKey)]
		}
		if strings.HasPrefix(srcObjectKey, o.TmpPath) {
			destObjectKey := strings.Replace(srcObjectKey, o.TmpPath, "", 1)
			//destObjectKey := srcObjectKey[len(o.TmpPath):len(srcObjectKey)]
			_, err := o.bucket.CopyObject(srcObjectKey, destObjectKey)
			if err != nil {
				return copyUrls, err
			}
			copyUrls = append(copyUrls, strings.Replace(u, o.TmpPath, "", 1))
		} else {
			copyUrls = append(copyUrls, u)
		}
	}
	return copyUrls, nil
}

// UploadBase64 Base64图片上传
// @return error
func (o *OSS) UploadBase64(data string, objectKey string) error {
	err := o.connection()
	decodeString, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	err = o.bucket.PutObject(objectKey, strings.NewReader(string(decodeString)))
	if err != nil {
		return err
	}
	return nil
}

// UploadString 上传字符串
func (o *OSS) UploadString(data string, objectKey string) error {
	err := o.connection()
	if err != nil {
		return err
	}
	err = o.bucket.PutObject(objectKey, strings.NewReader(data))
	if err != nil {
		return err
	}
	return nil
}

// UploadFile 上传文件
// objectKey
// filePath
// partSize 100K <= ? <= 5G
func (o *OSS) UploadFile(objectKey string, filePath string, partSize int64) error {
	err := o.connection()
	if err != nil {
		return err
	}
	err = o.bucket.UploadFile(objectKey, filePath, partSize)
	if err != nil {
		return err
	}
	return nil
}

// DownloadFile 下载文件到本地
func (o *OSS) DownloadFile(objectKey string, filePath string, partSize int64) error {
	err := o.connection()
	if err != nil {
		return err
	}
	err = o.bucket.DownloadFile(objectKey, filePath, partSize)
	if err != nil {
		return err
	}
	return nil
}

type ImageInfo struct {
	ExifTag struct {
		Value string `json:"value"`
	} `json:"ExifTag"`
	FileSize struct {
		Value string `json:"value"`
	} `json:"FileSize"`
	Format struct {
		Value string `json:"value"`
	} `json:"Format"`
	ImageHeight struct {
		Value string `json:"value"`
	} `json:"ImageHeight"`
	ImageWidth struct {
		Value string `json:"value"`
	} `json:"ImageWidth"`
	Orientation struct {
		Value string `json:"value"`
	} `json:"Orientation"`
}
