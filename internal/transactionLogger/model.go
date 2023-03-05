package transactionLogger

type EventType uint64

const (
	_                  = iota
	EventPut EventType = iota
	EventDelete
)

type Event struct {
	ID     uint64 `gorm:"primaryKey"`
	Method EventType
	Key    string
	Value  string
}
