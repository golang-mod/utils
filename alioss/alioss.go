package alioss

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

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
	onlineAt, err := time.ParseInLocation(time.RFC1123, date, time.Local)
	return onlineAt.Format("2006-01-02 15:04:05"), nil
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
		var uParse *url.URL
		uParse, err = url.Parse(u)
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
			_, err = o.bucket.CopyObject(srcObjectKey, destObjectKey)
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

// ReduceImg 缩放图片
func ReduceImg(imgUrl string) (string, error) {
	info, err := GetImageInfo(imgUrl)
	if err != nil {
		return "", err
	}
	width, err := strconv.ParseInt(info.ImageWidth.Value, 10, 64)

	if width > 1000 {
		imgUrl = imgUrl + "?x-oss-process=image/resize,w_1000,m_lfit"
	}
	return imgUrl, err
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

// GetImageInfo 获取图片信息
func GetImageInfo(imgUrl string) (info ImageInfo, err error) {

	res, err := resty.New().R().
		Get(imgUrl + "?x-oss-process=image/info")
	if err != nil {
		log.Println("err:", err)
	}
	if res.StatusCode() != http.StatusOK {
		err = http.ErrServerClosed
		return
	}

	err = json.Unmarshal([]byte(res.Body()), &info)
	if err != nil {
		log.Println("err:", err)
		return
	}
	return
}
