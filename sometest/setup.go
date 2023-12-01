package sometest

import "github.com/guatom999/Go-MicroService/config"

func NewTestConfig() *config.Config {
	cfg := config.LoadConfig("../env/test/.env")
	return &cfg
}
