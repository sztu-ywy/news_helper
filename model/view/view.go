package view

import (
	"github.com/uozi-tech/cosy/logger"
	"gorm.io/gorm"
)

// createViewFuncs functions list for creating views
var createViewFuncs = []func(db *gorm.DB) error{}

// CreateViews create views
func CreateViews(db *gorm.DB) {
	for _, f := range createViewFuncs {
		if err := f(db); err != nil {
			logger.Fatal(err)
		}
	}
	return
}
