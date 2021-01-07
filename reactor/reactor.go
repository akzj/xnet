package reactor

type Event struct {
	Fd int

	Read     bool
	Write    bool
	Callback Callback
}

type Reactor interface {
	AddEvent(event Event) error
	DelEvent(event Event) error

	Run() error
	Stop() error
}

type Callback interface {
	Read(bytes int64)
	Write(bytes int64)
	Close()
}
