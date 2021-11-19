package config

type Reader interface {
	Read() ([]string, error)
}

type Writer interface {
	Write([]string) error
}

// Just an example to show interface composition.
type ReaderWriter interface {
	Reader
	Writer
}
