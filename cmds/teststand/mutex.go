package main

import (
	"fmt"
	"sync"
)

// Mutex implements a high-level mutex with a few extra features to
// allow a better user interface.
//
// Extra features:
// - Order of locks maintained for next person.
// - Messsaging to waiting users trying to lock.
// - Only current holder can release the mutex.
type Mutex struct {
	queue  []*User
	holder string
	m      sync.Mutex
}

// User is someone waiting on the mutex.
type User struct {
	name string
	m    chan struct{}
	msgs chan string
}

// Updater is a function for receiving messages.
type Updater func(string) error

// Lock takes the mutex or waits until the mutex is available to lock.
func (m *Mutex) Lock(name string, update Updater) error {
	// Short path
	m.m.Lock()
	if m.holder == "" {
		m.holder = name
		m.m.Unlock()
		return nil
	}

	// Long path
	ch := make(chan struct{})
	msgs := make(chan string)
	m.queue = append(m.queue, &User{name, ch, msgs})
	m.m.Unlock()

	update(fmt.Sprintf("Held by %s", m.holder))
	m.Message(fmt.Sprintf("%s in now waiting", name))

	for {
		select {
		case <-ch:
			m.m.Lock()
			m.holder = name
			m.m.Unlock()
			m.Message(fmt.Sprintf("%s now has the mutex", name))
			return nil
		case msg := <-msgs:
			_ = update(msg) // TODO: Remove on error
		}
	}
}

// Unlock releases the mutex and allows the next user to take it.
func (m *Mutex) Unlock(name string) error {
	m.m.Lock()
	defer m.m.Unlock()

	if name != m.holder {
		return fmt.Errorf("can't unlock mutex: %s currently holds the key (you are %s)", m.holder, name)
	}
	m.Message(fmt.Sprintf("%s is done with the mutex", name))

	m.holder = ""
	if len(m.queue) > 0 {
		m.queue[0].m <- struct{}{}
		m.queue = m.queue[1:]
	}
	return nil
}

// Message sends messages to waiting users.
func (m *Mutex) Message(msg string) {
	for _, u := range m.queue {
		msgs := u.msgs
		go func() { msgs <- msg }()
	}
}
