package config

import (
	"testing"
	"baselib/lib-flag"
)

func TestLoadConfig(t *testing.T) {
	confPath := lib_flag.ConfPath
	LoadConfig(confPath)
}