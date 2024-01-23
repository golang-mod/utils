package alioss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"time"
)

// 用户上传文件时指定的前缀。
//var uploadDir = "user-dir-prefix/"
//var expireTime int64 = 30

func getGmtIso8601(expireEnd int64) string {
	var tokenExpire = time.Unix(expireEnd, 0).UTC().Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

// PolicyToken
// accessKeyId：设置您的AccessKeyId。
// accessKeySecret：设置您的AessKeySecret。
// host：格式为https://bucketname.endpoint，例如https://bucket-name.oss-cn-hangzhou.aliyuncs.com。关于Endpoint的介绍，请参见Endpoint访问域名。
// callbackUrl：设置上传回调URL，即回调服务器地址，用于处理应用服务器与OSS之间的通信。OSS会在文件上传完成后，把文件上传信息通过此回调URL发送给应用服务器。
// dir：若要设置上传到OSS文件的前缀则需要配置此项，否则置空即可。
type PolicyToken struct {
	AccessKeyId string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
	Callback    string `json:"callback"`
}

type CallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

func (o *OSS) getPolicyToken() PolicyToken {
	now := time.Now().Unix()
	expireEnd := now + o.Config.ExpireTime
	var tokenExpire = getGmtIso8601(expireEnd)

	//create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, o.Config.UploadDir)
	config.Conditions = append(config.Conditions, condition)

	//calucate signature
	result, err := json.Marshal(config)
	if err != nil {
		fmt.Println("calucate signature err:", err)
	}
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(o.Config.AccessKeySecret))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var policyToken PolicyToken
	policyToken.AccessKeyId = o.Config.AccessKeyId
	policyToken.Host = o.Config.EndpointUrl
	policyToken.Expire = expireEnd
	policyToken.Signature = signedStr
	policyToken.Directory = o.Config.UploadDir
	policyToken.Policy = debyte

	if o.Config.CallbackUrl != "" {
		var callbackParam CallbackParam
		callbackParam.CallbackUrl = o.Config.CallbackUrl
		// *重要，按Key递增排序规则排序，否则无法通过回调校验
		// filename=tmp%2Ftest%2FijIPTGn4NrWXc398adf3b96050b1914b61c709c48dcf.jpg&height=656&mimeType=image%2Fjpeg&size=73432&width=822
		callbackParam.CallbackBody = "filename=${object}&height=${imageInfo.height}&mimeType=${mimeType}&size=${size}&width=${imageInfo.width}"
		callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"
		callbackStr, err := json.Marshal(callbackParam)
		if err != nil {
			fmt.Println("callback json err:", err)
		}
		callbackBase64 := base64.StdEncoding.EncodeToString(callbackStr)
		policyToken.Callback = callbackBase64
	}

	//response,err:=json.Marshal(policyToken)
	//if err != nil {
	//    fmt.Println("json err:", err)
	//}
	return policyToken
}
