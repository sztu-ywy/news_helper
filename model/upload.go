package model

import (
	"github.com/uozi-tech/cosy/sonyflake"
	"gorm.io/gorm"
)

type Upload struct {
	Model
	UserId       uint64 `json:"user_id,omitempty" gorm:"index"`
	User         *User  `json:"user,omitempty"`
	MIME         string `json:"mime,omitempty" gorm:"index"`
	Name         string `json:"name,omitempty" gorm:"index"`
	Path         string `json:"path,omitempty" gorm:"index"`
	Thumbnail    string `json:"thumbnail,omitempty"`
	OriginalPath string `json:"original_path" gorm:"-"`
	Size         int64  `json:"size,omitempty"`
	To           string `json:"to,omitempty" gorm:"index"`
}

func (u *Upload) BeforeCreate(_ *gorm.DB) error {
	if u.ID == 0 {
		u.ID = sonyflake.NextID()
	}
	return nil
}

// func (u *Upload) AfterFind(tx *gorm.DB) error {
// 	u.OriginalPath = u.Path
// 	res, err := url.JoinPath(settings.OssSettings.BaseUrl, u.Path)
// 	if err != nil {
// 		return nil
// 	}
// 	u.Path = res
// 	if strings.Contains(u.MIME, "image") {
// 		// oss: svg+xml encode is not supported.
// 		if u.MIME == "image/svg+xml" {
// 			u.Thumbnail = u.Path
// 			return nil
// 		}
// 		// oss: Avif encode is not supported.
// 		if u.MIME != "image/avif" {
// 			u.Thumbnail = u.Path + "?x-oss-process=image/resize,w_256,h_256,m_fill,webp"
// 		} else {
// 			u.Thumbnail = u.Path + "?x-oss-process=style/webp"
// 		}
// 		u.Path += "?x-oss-process=style/webp"
// 	} else if strings.Contains(u.MIME, "video") {
// 		u.Thumbnail = u.Path + "?x-oss-process=video/snapshot,t_0,f_jpg,w_0,h_0,m_fast"
// 	}
// 	return nil
// }

func (u *Upload) AfterSave(tx *gorm.DB) (err error) {
	DropCache("upload", u.ID)
	return nil
}

func (u *Upload) AfterDelete(tx *gorm.DB) (err error) {
	DropCache("upload", u.ID)
	return nil
}
