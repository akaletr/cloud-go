package app

type App interface {
	Init() error

	Start() error
	Stop() error
}
