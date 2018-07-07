package log

import (
	"os"

	"github.com/Sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	// Logger expone el logrus.Logger:
	Logger    *logrus.Logger
	formatter *prefixed.TextFormatter
)

func init() {
	formatter = new(prefixed.TextFormatter)
	logrus.StandardLogger().Formatter = formatter
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	Logger = logrus.StandardLogger()
}
