package conf

import (
	_ "embed"
)

//go:embed app.dev.yaml
var AppDevYaml []byte

//go:embed app.prod.yaml
var AppProdYaml []byte
