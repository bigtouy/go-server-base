package configs

type App struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}
