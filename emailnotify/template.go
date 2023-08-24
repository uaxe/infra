package emailnotify

import "io"

// MessageImplementer
type MessageImplementer interface {
	Content(io.Writer) error
}
