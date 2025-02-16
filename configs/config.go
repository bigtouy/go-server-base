package configs

type ServerConfig struct {
	System System `mapstructure:"system"`
	Log    Log    `mapstructure:"log"`
	App    App    `mapstructure:"app"`
	Oss    Oss    `mapstructure:"oss"`
	S3     S3     `mapstructure:"s3"`
	Qwen   Qwen   `mapstructure:"qwen"`
	Sqlite Sqlite `mapstructure:"sqlite"`
}
