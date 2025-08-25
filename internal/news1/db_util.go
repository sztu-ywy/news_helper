package news1

import (
	"context"
	"fmt"
	"news_helper/model"
	"news_helper/query"
	"time"

	"github.com/spf13/cast"
	"github.com/uozi-tech/cosy"
	"github.com/uozi-tech/cosy/logger"
)

// getUserEmail 获取用户邮箱
func getUserInfo(userId string) (user *model.User, err error) {
	db := cosy.UseDB(context.Background())
	res := db.Model(&model.User{}).Where("id = ?", userId).Find(user)
	err = res.Error
	logger.Infof("获取用户邮箱: UserID=%s, user=%s", userId, user)
	return
}

// getMediaInfo 爬取指定媒体的新闻
func getMediaInfo(mediaID uint64) (model.News, error) {
	logger.Infof("开始爬取媒体 %d 的新闻", mediaID)
	db := cosy.UseDB(context.Background())
	news := model.News{}
	db.Model(&model.News{}).Where("media_id = ?", mediaID).Find(&news)

	return news, nil
}

// saveAnalysisResult 保存分析结果到数据库
func saveAnalysisResult(userId, taskId, analysisResult string, newsData []*model.News) error {
	logger.Infof("保存分析结果: UserID=%s, TaskID=%s", userId, taskId)

	// 创建分析结果记录
	analysisRecord := model.AnalysisResult{
		UserID:       cast.ToUint64(userId),
		TaskID:       cast.ToUint64(taskId),
		Title:        fmt.Sprintf("新闻分析报告 - %s", time.Now().Format("2006-01-02 15:04:05")),
		Content:      analysisResult,
		NewsCount:    len(newsData),
		AnalysisDate: uint64(time.Now().Unix()),
		Status:       "completed",
	}

	db := cosy.UseDB(context.Background())
	db.Model(&model.AnalysisResult{}).Create(&analysisRecord)

	// 保存新闻数据到数据库
	query.News.Create(newsData...)

	logger.Infof("分析结果保存完成: ID=%d", analysisRecord.ID)
	return nil
}
