package query

import (
	"errors"

	"news_helper/internal/acl"
	"news_helper/model"
	"github.com/uozi-tech/cosy/logger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Init(db *gorm.DB) {
	SetDefault(db)

	ug := UserGroup

	// Creating an initial admin group
	_, err := ug.Unscoped().Where(ug.ID.Eq(1)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Info("Creating an initial admin group")
		err = ug.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"permissions"}),
		}).Create(&model.UserGroup{
			Model: model.Model{
				ID: 1,
			},
			Name: "管理员",
			Permissions: acl.PermissionList{
				acl.Permission{
					Subject: acl.All,
					Action:  acl.Write,
				},
			},
		})
		if err != nil {
			logger.Fatal(err)
		}
	}
	// Create initial user
	u := User
	initUser, err := u.Unscoped().Where(u.ID.Eq(1)).First()
	if err == nil {
		_ = db.Model(initUser).Updates(&model.User{
			UserGroupID: 1,
		})
	}

	// TODO   这里修改admin的原始信息
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Info("Creating initial user, email is admin, password is admin")
		pwd, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			logger.Fatal(err)
		}

		_, err = u.Unscoped().Where(u.ID.Eq(1)).Assign(field.Attrs(&model.User{
			UserGroupID: 1,
		})).Attrs(field.Attrs(&model.User{
			// todo: 替换为真实的超管邮箱，可以在配置文件配置，以便后续发送邮件通知
			Email:    "admin",
			Password: string(pwd),
			Name:     "admin",
			// HotelID:   1,
			Status: model.UserStatusActive,
		})).FirstOrCreate()

		if err != nil {
			logger.Fatal("Fail to create initial user", err)
		}
	}

	InitHotel(db)
}

func InitHotel(db *gorm.DB) {
	// h := Hotel
	// adminHotel, err := h.Unscoped().Where(h.ID.Eq(1)).First()
	// if err == nil {
	// 	_ = db.Model(adminHotel).Updates(&model.Hotel{
	// 		Name:     "admin",
	// 		Location: "admin",
	// 	})
	// }

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	logger.Info("Creating initial admin hotel, name is admin, location is admin")
	// 	err = h.Clauses(clause.OnConflict{
	// 		Columns:   []clause.Column{{Name: "id"}},
	// 		DoUpdates: clause.AssignmentColumns([]string{"name", "location"}),
	// 	}).Create(&model.Hotel{
	// 		Model: model.Model{
	// 			ID: 1,
	// 		},
	// 		Name:     "admin",
	// 		Location: "admin",
	// 	})
	// 	if err != nil {
	// 		logger.Fatal("Fail to create initial admin hotel", err)
	// 	}
	// hotel := &model.Hotel{
	// 	Name:     "admin",
	// 	Location: "admin",
	// }
	// // 假设有一个方法来保存酒店信息到数据库
	// result := db.Create(hotel)
	// if result.Error != nil {
	// 	logger.Fatal("Failed to create hotel:", result.Error)
	// }
	// }
}
