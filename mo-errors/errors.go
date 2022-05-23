package mo_errors

import "errors"

//type NoResultError struct{}
//
//func (e *NoResultError) Error() string {
//	return "No Result"
//}
//
//type TimeoutError struct{}
//
//func (e *TimeoutError) Error() string {
//	return "Timeout"
//}
//
//type FullError struct{}
//
//func (e *FullError) Error() string {
//	return "Full Error"
//}

var (
	NoResultError = errors.New("no Result")
	TimeoutError  = errors.New("timeout")
	FullError     = errors.New("full error")
)
