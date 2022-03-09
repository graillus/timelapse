package log

import (
	"github.com/sirupsen/logrus"
)

//nolint:gochecknoinits
func init() {
	logger = logrus.StandardLogger()
}
