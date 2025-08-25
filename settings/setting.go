package settings

import (
	"strings"
	"time"

	"git.uozi.org/uozi/crypto"
	"github.com/spf13/cast"
	"github.com/uozi-tech/cosy/settings"
)

var (
	buildTime    string
	LastModified string
)

func init() {
	t := time.Unix(cast.ToInt64(buildTime), 0)
	LastModified = strings.ReplaceAll(t.Format(time.RFC1123), "UTC", "GMT")

	settings.Register("crypto", crypto.Settings)
	// settings.Register("sls", SLSSettings)
	settings.Register("email", EmailSettings)
	settings.Register("qwen", QwenSettings)
}
