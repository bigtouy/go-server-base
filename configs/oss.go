package configs

type Oss struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	Region          string `mapstructure:"region"`
	BucketName      string `mapstructure:"bucket_name"`
	ScType          string `mapstructure:"sc_type"`
}
