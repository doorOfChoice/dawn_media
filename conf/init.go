package conf

import (
	"github.com/BurntSushi/toml"
)

type Conf struct {
	Address   string `toml:"address"`
	UploadDir string `toml:"upload_dir"`
}

var conf Conf

func Init() {
	if _, err := toml.DecodeFile("media.toml", &conf); err != nil {
		panic(err)
	}
}

func C() Conf {
	return conf
}
