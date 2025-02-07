package jsonhelper

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

func Encode[T any](t T) []byte {
	b, err := json.Marshal(t)
	if err != nil {
		logrus.WithError(err).WithField("t", t).Fatalln("couldn't encode the variable")

		return nil
	}

	return b
}

func Decode[T any](b []byte) T {
	var t T

	if err := json.Unmarshal(b, &t); err != nil {
		logrus.WithError(err).WithField("t", t).Fatalln("couldn't decode the variable")
	}

	return t
}
