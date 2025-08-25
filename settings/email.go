package settings

type Email struct {
	Host     string `ini:"host"`
	Port     string `ini:"port"`
	Email    string `ini:"email"`
	Password string `ini:"password"`
}

var EmailSettings = &Email{}
