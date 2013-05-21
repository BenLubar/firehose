// Package firehose provides a stream of atom feed entries.
package firehose

import "code.google.com/p/go.blog/pkg/atom"

type Feed interface {
	// A channel where the entries will be sent. It is safe to cache the
	// return value of this method. The channel will be closed when the
	// Close method is called.
	Entries() <-chan *atom.Entry

	// Closes the Feed and returns the last error encountered, if any.
	Close() error
}
