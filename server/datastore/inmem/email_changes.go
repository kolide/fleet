package inmem

func (ds *Datastore) PendingEmailChange(uid uint, newEmail, token string) error {
	panic("deprecated")
}

func (ds *Datastore) ChangeUserEmail(token string) (string, error) {
	panic("deprecated")
}
