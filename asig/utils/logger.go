package utils

import(
	"github.com/sirupsen/logrus"
	"os"
	"fmt"
)

var file *os.File

func InitiateLogger(filePath string) (err error) {
	if file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		fmt.Println(err)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(file)
	return err
}

func CloseLogger() error {
	return file.Close()
}

func Log(status, code, message string, detail map[string]interface{}) {
	if status == "success" {
		logrus.WithFields(logrus.Fields{
			"status":          status,
			"code":            code,
			"detail":		   detail,
		}).Info(message)
	} else {
		logrus.WithFields(logrus.Fields{
			"status":          status,
			"code":            code,
			"detail":		   detail,
		}).Error(message)
	}
}