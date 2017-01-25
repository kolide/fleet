package inmem

// Stubs: inmem is going away
func (ds *Datastore) PendingEmailChange(uid uint, newEmail, token string) error {
	return nil
}

func (ds *Datastore) CommitEmailChange(token string) (string, error) {
	return "", nil
}
