package model

type httpMethod uint8

const (
	GET httpMethod = iota
	POST
	PUT
	DELETE
)

// Media  媒体
type Media struct {
	Model
	Name       string     `json:"name" gorm:"column:name;type:varchar(100);not null;comment:媒体名称"`
	Url        string     `json:"url" gorm:"column:url;type:varchar(100);not null;comment:媒体链接"`
	HttpMethod httpMethod `json:"http_method" gorm:"column:http_method;type:tinyint(1);not null;comment:请求方式"`
}
