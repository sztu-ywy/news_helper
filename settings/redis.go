package settings

type Redis struct {
	Host     string `json:"host"`
	Password string `json:"-,omitempty"`
	Port     uint   `json:"port"`
	User     string `json:"user"`
	DB       uint   `json:"db"`
	Prefix   string `json:"prefix"`
}

var RedisSettings = &Redis{}
