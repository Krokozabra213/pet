package chatusecases

func (a *Chat) Unsubscribe(uuid uint64) {
	a.clientRepo.Unsubscribe(uuid)
}
