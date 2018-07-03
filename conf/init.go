package conf

import (
	"github.com/BurntSushi/toml"
)

type Conf struct {
	Address     string `toml:"address"`
	Username    string `toml:"username"`
	Password    string `toml:"password"`
	Database    string `toml:"database"`
	PassSalt    string `toml:"pass_salt"`
	LoginSalt   string `toml:"login_salt"`
	SessionName string `toml:"session_name"`
	AvatarDir   string `toml:"avatar_dir"`
	AvatarMap   string `toml:"avatar_map"`
	CoverDir    string `toml:"cover_dir"`
	CoverMap    string `toml:"cover_map"`
	MediaDir    string `toml:"media_dir"`
	MediaMap    string `toml:"media_map"`
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
