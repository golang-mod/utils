package configs

type AliossConfig struct {
	AccessKeyId     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	Bucket          string `yaml:"bucket"`
	Endpoint        string `yaml:"endpoint"`
	EndpointUrl     string `yaml:"endpoint_url"`
	CallbackUrl     string `yaml:"callback_url"`
}
