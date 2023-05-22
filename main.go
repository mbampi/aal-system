package main

import "github.com/sirupsen/logrus"

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.Info("Starting AAL System")
}
