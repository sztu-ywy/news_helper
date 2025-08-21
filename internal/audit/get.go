package audit

import (
	"context"

	"news_helper/internal/geoip"

	"news_helper/model"
	"news_helper/settings"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	cModel "github.com/uozi-tech/cosy/model"
)

func GetLogs(ctx context.Context, from int64, to int64, offset, pageSize int64, queryExp string) (resp *sls.GetLogsResponse, err error) {
	endpoint := settings.SLSSettings.EndPoint
	projectName := settings.SLSSettings.ProjectName
	logStoreName := settings.SLSSettings.LogStoreName

	provider := getCredentialsProvider()

	client := sls.CreateNormalInterfaceV2(endpoint, provider)

	resp, err = client.GetLogs(projectName, logStoreName, Topic, from, to,
		queryExp, pageSize, offset, true)

	if err != nil {
		return
	}

	// collect user ids
	var userIDs []string
	for _, log := range resp.Logs {
		userIDs = append(userIDs, log["user_id"])
	}

	if len(userIDs) == 0 {
		return
	}

	userIDs = lo.Uniq(userIDs)

	var users []*model.User
	db := cModel.UseDB(ctx)
	db.Where("id IN (?)", userIDs).Find(&users)

	userIdMap := lo.SliceToMap(users, func(item *model.User) (string, *model.User) {
		return cast.ToString(item.ID), item
	})

	for _, log := range resp.Logs {
		// username
		if user, ok := userIdMap[log["user_id"]]; ok {
			log["user"] = user.Name
		} else {
			log["user"] = ""
		}

		// geoip
		log["geoip"] = geoip.ParseIP(log["ip"])
	}

	return
}
