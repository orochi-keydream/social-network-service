package model

import "fmt"

// TODO: How to properly wrap these erorrs and make it unwrapable?

type ClientError struct {
	Message string
	Err     error
}

func NewClientError(msg string, err error) *ClientError {
	return &ClientError{
		Message: msg,
		Err:     err,
	}
}

func (e *ClientError) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return fmt.Sprintf("%v: %v", e.Message, e.Err)
}

type NotFoundError struct {
	Message string
	Err     error
}

func NewNotFoundError(msg string, err error) *NotFoundError {
	return &NotFoundError{
		Message: msg,
		Err:     err,
	}
}

func (e *NotFoundError) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return fmt.Sprintf("%v: %v", e.Message, e.Err)
}

type ForbiddenError struct {
	Message string
	Err     error
}

func NewForbiddenError(msg string, err error) *ForbiddenError {
	return &ForbiddenError{
		Message: msg,
		Err:     err,
	}
}

func (e *ForbiddenError) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return fmt.Sprintf("%v: %v", e.Message, e.Err)
}

type UnauthenticatedError struct {
	Message string
	Err     error
}

func NewUnauthenticatedError(msg string, err error) *UnauthenticatedError {
	return &UnauthenticatedError{
		Message: msg,
		Err:     err,
	}
}

func (e *UnauthenticatedError) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return fmt.Sprintf("%v: %v", e.Message, e.Err)
}
