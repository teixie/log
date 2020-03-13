package log

import (
	"bytes"
	"fmt"
	"time"

	"github.com/rs/xid"
)

type bufferLogger struct {
	buf     *bytes.Buffer
	logId   string
	startAt time.Time
	timers  map[string]time.Time
}

func (b *bufferLogger) StartTimer(name string) {
	b.timers[name] = time.Now()
}

func (b *bufferLogger) StopTimer(name string) {
	if timer, ok := b.timers[name]; ok {
		duration := float64(time.Now().Sub(timer).Nanoseconds()) / 1e6
		b.buf.WriteString(fmt.Sprintf(" %s=%vms", name, duration))
	}
}

func (b *bufferLogger) Append(args ...interface{}) {
	b.buf.WriteString(" " + appendString(args...))
}

func (b *bufferLogger) AppendAll(args ...interface{}) {
	for _, v := range args {
		b.buf.WriteString(fmt.Sprintf(" %+v", v))
	}
}

func (b *bufferLogger) Flush() {
	duration := float64(time.Now().Sub(b.startAt).Nanoseconds()) / 1e6
	Log.Info(fmt.Sprintf("logId=%s cost=%vms %s", b.logId, duration, b.buf.String()))
}

func (b *bufferLogger) Debug(format string, args ...interface{}) {
	Log.Debugf(fmt.Sprintf("logId=%s %s", b.logId, format), args...)
}

func (b *bufferLogger) Info(format string, args ...interface{}) {
	Log.Infof(fmt.Sprintf("logId=%s %s", b.logId, format), args...)
}

func (b *bufferLogger) Notice(format string, args ...interface{}) {
	Log.Noticef(fmt.Sprintf("logId=%s %s", b.logId, format), args...)
}

func (b *bufferLogger) Warning(format string, args ...interface{}) {
	Log.Warningf(fmt.Sprintf("logId=%s %s", b.logId, format), args...)
}

func (b *bufferLogger) Error(format string, args ...interface{}) {
	Log.Errorf(fmt.Sprintf("logId=%s %s", b.logId, format), args...)
}

func (b *bufferLogger) Critical(format string, args ...interface{}) {
	Log.Criticalf(fmt.Sprintf("logId=%s %s", b.logId, format), args...)
}

func appendString(args ...interface{}) string {
	if len(args) == 2 {
		return fmt.Sprintf("%+v=%+v", args[0], args[1])
	} else if len(args) == 1 {
		return fmt.Sprintf("%+v", args[0])
	}
	return ""
}

func NewBufferLogger(args ...interface{}) *bufferLogger {
	return &bufferLogger{
		buf:     bytes.NewBufferString(appendString(args...)),
		logId:   xid.New().String(),
		startAt: time.Now(),
		timers:  make(map[string]time.Time),
	}
}
