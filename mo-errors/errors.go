package mo_errors

type NoResultError struct{}

func (e NoResultError) Error() string {
	return "No Result"
}
