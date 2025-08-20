package acl

func Can(m Map, subject Subject, action Action) bool {
	// If the subject has the manage action, then the user can do create, read, update, delete within the scope
	manage := m[subject][Write] || m[All][Write]

	if manage {
		return true
	}

	return m[subject][action] || m[All][action]
}
