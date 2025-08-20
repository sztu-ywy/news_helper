package acl

func (p *PermissionList) ToMap() Map {
	m := make(Map)
	for _, v := range *p {
		if _, ok := m[v.Subject]; !ok {
			m[v.Subject] = make(map[Action]bool)
		}
		m[v.Subject][v.Action] = true
	}
	return m
}
