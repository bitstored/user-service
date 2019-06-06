package repository

import (
	"log"
)

type Kind string
type Status string

const (
	EVENT_TYPE_INSERT Kind   = "EVENT_TYPE_INSERT"
	EVENT_TYPE_UPDATE Kind   = "EVENT_TYPE_UPDATE"
	EVENT_TYPE_DELETE Kind   = "EVENT_TYPE_INSERT"
	EVENT_TYPE_GET    Kind   = "EVENT_TYPE_GET"
	STATUS_SUCCESS    Status = "SUCCESS"
	STATUS_ERROR      Status = "ERROR"
)

func LogEvent(kind Kind, status Status, err error, entity interface{}) {
	log.Printf("%s: %v\n\tstatus:%s\n\terror:%v\n", kind, entity, status, err)
}
