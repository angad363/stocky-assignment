package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
    Log = logrus.New()

    // Output to both file and stdout
    file, err := os.OpenFile("stocky.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		file = os.Stdout
	}

	Log.SetOutput(io.MultiWriter(os.Stdout, file))

    Log.SetLevel(logrus.InfoLevel)
    Log.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
        ForceColors:   true,
        TimestampFormat: "2006-01-02 15:04:05",
    })
}