package repository

var RepoErr = NewRepoError("database err")

type RepoError struct {
	message string
}

func (e *RepoError) Error() string {
	return e.message
}

func NewRepoError(msg string) *RepoError {
	return &RepoError{message: msg}
}
