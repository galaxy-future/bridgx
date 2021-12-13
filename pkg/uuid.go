package pkg

import (
	"errors"

	"github.com/google/uuid"
	uuid4 "github.com/satori/go.uuid"
)

var (
	uuidCollector        = map[string]struct{}{}
	errGenerateUUIDPanic = errors.New("generate uuid failed with panic")
)

// NewUUID use google/uuid UUID version 4.
func NewUUID() (uid string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverPanic2Err(r)
		}
	}()
	uid = uuid.NewString()
	if isCollected(uid) {
		// TODO: backoff retry?
	}
	collectUUID(uid)
	return uid, nil
}

// NewUUID4 use satori/uuid UUID version 4.
func NewUUID4() (uid string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverPanic2Err(r)
		}
	}()
	uid = uuid4.NewV4().String()
	return uid, nil
}

// collectUUID is not thread safe!
func collectUUID(uid ...string) {
	for _, id := range uid {
		if !isCollected(id) {
			uuidCollector[id] = struct{}{}
		}
	}
}

// isCollected is not thread safe!
func isCollected(uid string) bool {
	_, ok := uuidCollector[uid]
	return ok
}

func recoverPanic2Err(r interface{}) error {
	err, ok := r.(error)
	if ok {
		return err
	}
	return errGenerateUUIDPanic
}
