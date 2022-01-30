package main

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func TimestampToDate(ts int64) string {
	return time.Unix(ts, 0).UTC().String()
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panic(msg)
	}
}
