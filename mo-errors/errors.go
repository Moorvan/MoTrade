package mo_errors

type NoResultError struct{}

func (e NoResultError) Error() string {
	return "No Result"
}

type TimeoutError struct{}

func (e TimeoutError) Error() string {
	return "Timeout"
}
