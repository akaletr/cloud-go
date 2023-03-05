package transactionLogger

type TransactionLogger interface {
	WritePut(key, value string)
	WriteDelete(key string)

	Err() <-chan error
	Run()
	ReadEvents() (<-chan Event, <-chan error)
}
