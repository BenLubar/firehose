package firehose

import (
	"sync"

	"code.google.com/p/go.blog/pkg/atom"
)

type multi struct {
	feeds   []Feed
	entries chan *atom.Entry
	done    sync.WaitGroup
	close   chan struct{}
}

// Multi joins zero or more feeds together. Closing a multi-feed will close
// all of the feeds given as arguments to this function.
func Multi(feeds ...Feed) Feed {
	feed := &multi{
		feeds:   feeds,
		entries: make(chan *atom.Entry),
		close:   make(chan struct{}),
	}
	feed.done.Add(len(feeds))
	for _, f := range feeds {
		go feed.fetch(f)
	}
	go feed.cleanup()
	return feed
}

func (m *multi) fetch(f Feed) {
	defer m.done.Done()

	entries := f.Entries()
	var queue *atom.Entry

	for {
		var fillQueue <-chan *atom.Entry
		var emptyQueue chan<- *atom.Entry
		if queue == nil {
			fillQueue = entries
		} else {
			emptyQueue = m.entries
		}

		select {
		case queue = <-fillQueue:
		case emptyQueue <- queue:
			queue = nil

		case <-m.close:
			return
		}
	}
}

func (m *multi) cleanup() {
	m.done.Wait()
	close(m.entries)
}

func (m *multi) Entries() <-chan *atom.Entry {
	return m.entries
}

func (m *multi) Close() error {
	close(m.close)

	var lastError error

	for _, f := range m.feeds {
		if err := f.Close(); err != nil {
			lastError = err
		}
	}

	return lastError
}
