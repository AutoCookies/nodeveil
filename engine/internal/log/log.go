package log

import (
	"io"
	stdlog "log"
)

func New(out io.Writer) *stdlog.Logger {
	return stdlog.New(out, "nodeveil-engine ", stdlog.LstdFlags|stdlog.LUTC)
}
