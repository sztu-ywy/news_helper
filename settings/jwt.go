package settings

type Jwt struct {
	Issuer string `ini:"ywy"`
	Secret string `ini:"secret"`
	Expire int64  `ini:"expire"`
}

var JwtSettings = &Jwt{}
