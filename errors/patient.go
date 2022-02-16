package errors

import (
	errs "errors"
)

type ServiceError error

var (
	ErrNotFound              ServiceError = errs.New("patient not found")
	ErrEmptyData             ServiceError = errs.New("data is empty")
	ErrDataNotInitialized    ServiceError = errs.New("data not initialized")
	ErrAlreadyExists         ServiceError = errs.New("patient already exists")
	ErrCreationFailed        ServiceError = errs.New("error creating patient")
)
