package labels

import (
	"bufio"
)

type MyWriteCloser struct {
	*bufio.Writer
}

func (mwc *MyWriteCloser) Close() error {
	// Noop
	return nil
}
