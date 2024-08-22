package log

import (
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/Kajmany/cata-up/cfg"
)

// What's the point: it's an slogger to pass around, but there's also a string bfufer
// for my UI to pull from.

const maxSize int = 100

// stringBuffer isn't a *real* buffer, more just a string-wrangler.
type stringBuffer struct {
	mu     *sync.Mutex
	buffer []string
	limit  int
}

func newStringBuffer(limit int) *stringBuffer {
	return &stringBuffer{
		mu:     &sync.Mutex{},
		buffer: make([]string, 0, limit),
		limit:  limit,
	}
}

func (sb *stringBuffer) Write(p []byte) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	str := string(p)
	if len(sb.buffer) == sb.limit {
		sb.buffer = sb.buffer[1:]
	}
	sb.buffer = append(sb.buffer, str)
	return len(p), nil
}

func (sb *stringBuffer) GetBuffer() []string {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return append([]string(nil), sb.buffer...)
}

type BufferedLogger struct {
	L *slog.Logger
	B *stringBuffer
}

func NewBufferedLogger(config *cfg.Config) (BufferedLogger, error) {
	var out io.Writer
	buffer := newStringBuffer(maxSize)
	// TODO: I hate the way this struct is addressed
	switch config.LogExport.LogOutPutMode {

	case cfg.Log2Path:
		file, err := os.OpenFile(config.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return BufferedLogger{}, err
		}
		out = io.MultiWriter(file, buffer)
	case cfg.Log2Stderr:
		out = io.MultiWriter(os.Stderr, buffer)
	case cfg.LogOff:
		out = buffer
	}

	opts := slog.HandlerOptions{Level: config.LogLevel.Level}
	handler := slog.NewTextHandler(out, &opts)
	logger := slog.New(handler)
	return BufferedLogger{logger, buffer}, nil
}
