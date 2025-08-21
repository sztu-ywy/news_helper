package audit

import (
	"git.uozi.org/uozi/burn-api/settings"
	sls "github.com/aliyun/aliyun-log-go-sdk"
)

func GetHistograms(from int64, to int64, queryExp string) (resp *sls.GetHistogramsResponse, err error) {
	endpoint := settings.SLSSettings.EndPoint
	projectName := settings.SLSSettings.ProjectName
	logStoreName := settings.SLSSettings.LogStoreName

	provider := getCredentialsProvider()
	client := sls.CreateNormalInterfaceV2(endpoint, provider)

	resp, err = client.GetHistograms(projectName, logStoreName, Topic, from, to, queryExp)

	return
}
