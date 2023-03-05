package transactionLogger

type postgres struct {
}

func NewPostgres() TransactionLogger {
	return &postgres{}
}

func (p postgres) WritePut(key, value string) {
	//TODO implement me
	panic("implement me")
}

func (p postgres) WriteDelete(key string) {
	//TODO implement me
	panic("implement me")
}

func (p postgres) Err() <-chan error {
	//TODO implement me
	panic("implement me")
}

func (p postgres) Run() {
	//TODO implement me
	panic("implement me")
}

func (p postgres) ReadEvents() (<-chan Event, <-chan error) {
	//TODO implement me
	panic("implement me")
}
