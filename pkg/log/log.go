package log

import (
	"bufio"
	"io"
	"log"
	"log/slog"
)

func init() {
	pr, pw := io.Pipe()
	log.SetOutput(pw)
	go convertLogToSlog(pr)
}

func convertLogToSlog(r io.Reader) {
	logger := slog.With("plain-log", true)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		logger.Info(sc.Text())
	}
	panic(sc.Err())
}

func Err(err error) slog.Attr {
	return slog.Any("error", err)
}
