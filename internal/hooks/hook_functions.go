package hooks

import (
	"fmt"
	"fp-dim-aws-guard-duty-ingress/internal"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type LogEvent struct {
	ModuleName string    `json:"module_name"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	Caller     string    `json:"caller"`
	Time       time.Time `json:"time"`
}

// hook to buffer logs and only send at right severity.
type LoggingHook struct {
}

// Fire will append all logs to a circular buffer and only 'flush'
// them when a log of sufficient severity(ERROR) is emitted.
func (h *LoggingHook) Fire(entry *logrus.Entry) error {
	go func() {
		log := &LogEvent{
			ModuleName: os.Getenv("MODULE_SVC_NAME"),
			Level:      entry.Level.String(),
			Message:    entry.Message,
			Caller:     entry.Caller.Func.Name(),
			Time:       entry.Time,
		}

		_, err := internal.MakeRequest("POST", "logevent", log)

		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	return nil
}

// Levels define on which log levels this LoggingHook would trigger
func (h *LoggingHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel}
}
