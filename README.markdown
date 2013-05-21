# firehose
--
    import "github.com/BenLubar/firehose"

Package firehose provides a stream of atom feed entries.

## Usage

#### type Feed

```go
type Feed interface {
	// A channel where the entries will be sent. It is safe to cache the
	// return value of this method. The channel will be closed when the
	// Close method is called.
	Entries() <-chan *atom.Entry

	// Closes the Feed and returns the last error encountered, if any.
	Close() error
}
```


#### func  Fetch

```go
func Fetch(url string, delay time.Duration) Feed
```
Fetch downloads an atom feed over HTTP at the specified interval and sends
entries that have a unique ID to the Entries channel.

#### func  Multi

```go
func Multi(feeds ...Feed) Feed
```
Multi joins zero or more feeds together. Closing a multi-feed will close all of
the feeds given as arguments to this function.

--
**godocdown** http://github.com/robertkrimen/godocdown
