package repository

// RepositoryResult is a struct to wrap repository result
// so its easy to use it in channel
type RepositoryResult struct {
	Result interface{}
	Error  error
}
