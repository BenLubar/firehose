package firehose

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"code.google.com/p/go.blog/pkg/atom"
)

type fetch struct {
	entries chan *atom.Entry
	close   chan chan error
}

// Fetch downloads an atom feed over HTTP at the specified interval and sends
// entries that have a unique ID to the Entries channel.
func Fetch(url string, delay time.Duration) Feed {
	f := &fetch{
		entries: make(chan *atom.Entry),
		close:   make(chan chan error),
	}
	go f.fetch(url, delay)
	return f
}

func (f *fetch) fetch(url string, delay time.Duration) {
	defer close(f.close)
	defer close(f.entries)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/atom+xml")

	seen := make(map[string]bool)

	var lastErr error
	fetchAfter := time.After(0) // first fetch is immediate

	var queue []*atom.Entry

	for {
		var queued *atom.Entry
		var sendEntry chan *atom.Entry
		if len(queue) > 0 {
			queued = queue[0]
			sendEntry = f.entries
		}

		select {
		case c := <-f.close:
			c <- lastErr
			return

		case sendEntry <- queued:
			queue = queue[1:]

		case <-fetchAfter:
			fetchAfter = time.After(delay)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				lastErr = err
				// TODO: back off
				continue
			}

			if resp.StatusCode != http.StatusOK {
				if resp.StatusCode != http.StatusNotModified {
					log.Printf("Fetch %s: %d %s", url, resp.StatusCode, resp.Status)
				}
				io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
				continue
			}

			req.Header.Set("If-None-Match", resp.Header.Get("ETag"))
			var feed atom.Feed
			err = xml.NewDecoder(resp.Body).Decode(&feed)
			resp.Body.Close()

			for _, entry := range feed.Entry {
				if !seen[entry.ID] {
					queue = append(queue, entry)
					seen[entry.ID] = true
				}
			}
		}
	}
}

func (f *fetch) Close() error {
	ch := make(chan error)
	f.close <- ch
	return <-ch
}

func (f *fetch) Entries() <-chan *atom.Entry {
	return f.entries
}
