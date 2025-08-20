package acl

// Permission is the smallest unit of access control
type Permission struct {
	Subject Subject `json:"subject"`
	Action  Action  `json:"action"`
}

type PermissionList []Permission

type Map map[Subject]map[Action]bool
