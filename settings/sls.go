package settings

type SLS struct {
	AccessKeyId     string
	AccessKeySecret string
	EndPoint        string
	ProjectName     string
	LogStoreName    string
	Source          string
}

var SLSSettings = &SLS{}

func (s *SLS) Enable() bool {
	return s.AccessKeyId != "" && s.AccessKeySecret != "" && s.EndPoint != "" && s.ProjectName != ""
}
