package configs

type Qwen struct {
	Region      string `mapstructure:"region"`
	NlsAppKey   string `mapstructure:"nls_app_key"`
	ApiKey      string `mapstructure:"api_key"`
	WorkspaceId string `mapstructure:"workspace_id"`
}
