// settings/server.go
package settings

type Server struct {
	RunMode  string
	Host     string
	HttpPort string
}

var ServerSetting = &Server{}
