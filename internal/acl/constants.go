package acl

type Subject string

type Action string

const (
	All              Subject = "all"
	ServerAnalytics  Subject = "server_analytics"
	User             Subject = "user"
	Media            Subject = "media"
	Device           Subject = "device"
	RoleTemplate     Subject = "role_template"
	IntelligentAgent Subject = "intelligent_agent"
	OTA              Subject = "ota"
	PostAndPage      Subject = "post_and_page"
	AuditLog         Subject = "audit_log"
)

const (
	Read  Action = "read"
	Write Action = "write"
)

var subjects = []Subject{
	All,
	ServerAnalytics,
	User,
	Media,
	Device,
	RoleTemplate,
	IntelligentAgent,
	OTA,
	PostAndPage,
	AuditLog,
}

func GetSubjects() []Subject {
	return subjects[:]
}
