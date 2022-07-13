package tcp

import (
	"io"
	"strings"
)

func isServerClose(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}

func isClientClose(err error) bool {
	return err == io.EOF
}
