package internal

import "github.com/f1bonacc1/glippy"

type Clipboard interface {
	// write to the system clipboard
	Write(data []byte) error
}

// a mock clipboard interface until i implement properly
type clipper struct{}

func (c *clipper) Write(data []byte) error {
	s := string(data)
	return glippy.Set(s)
}
