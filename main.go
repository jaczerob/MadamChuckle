package main

import (
	"github.com/jaczerob/madamchuckle/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("error handling command")
	}
}
