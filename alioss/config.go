package alioss

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

type StsConfig struct {
	RegionId        string `yaml:"region_id"`
	AccessKeyId     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	RoleArn         string `yaml:"role_arn"`
	RoleSessionName string `yaml:"role_session_name"` // 自定义角色会话名称，用来区分不同的令牌，例如可填写为SessionTest。
}
