package configs

type Sqlite struct {
	DbPath string `mapstructure:"db_path"`
	DbName string `mapstructure:"db_name"`
}
