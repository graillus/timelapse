package log

import (
	"github.com/sirupsen/logrus"
)

func init() {
	logger = logrus.StandardLogger()
}
