package reactor

import (
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"
)

type Kqueue struct {
	handle int
	events []syscall.Kevent_t

	wait   int //read event channel
	notify int //write event channel

	notified int32 //notify status

	eventLocker sync.Mutex
	change      []syscall.Kevent_t
}

func (k *Kqueue) AddEvent(event Event) error {
	k.eventLocker.Lock()
	kevent := syscall.Kevent_t{
		Ident:  uint64(event.Fd),
		Udata:  &event,
	}
	k.change = append(k.change, )
	k.eventLocker.Unlock()

	if atomic.CompareAndSwapInt32(&k.notified, 0, 1) {
		_, err := syscall.Write(k.notify, []byte{'K'})
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Kqueue) DelEvent(event Event) error {
	panic("implement me")
}

func (k *Kqueue) Run() error {
	for {
		n, err := syscall.Kevent(k.handle, nil, k.events, nil)
		if err != nil {
			if err == syscall.EINTR {
				continue
			}
			panic(err)
		}
		for i := 0; i < n; i++ {
			event := k.events[i]

			ident := int(event.Ident)
			if ident == k.wait {
				//todo handle event
				continue
			}
			callback := *(*Callback)(unsafe.Pointer(&event.Udata))
			switch event.Fflags {
			case syscall.EVFILT_READ:
				callback.Read(event.Data)
			case syscall.EVFILT_WRITE:
				callback.Write(event.Data)
			}
		}
	}
}

func (k *Kqueue) Stop() error {
	panic("implement me")
}

func New() (Reactor, error) {
	handle, err := syscall.Kqueue()
	if err != nil {
		return nil, err
	}
	var channel [2]int
	if err := syscall.Pipe(channel[:]); err != nil {
		return nil, err
	}
	return &Kqueue{
		handle: handle,
		events: make([]syscall.Kevent_t, 1024),
		wait:   channel[0],
		notify: channel[1],
	}, nil
}
