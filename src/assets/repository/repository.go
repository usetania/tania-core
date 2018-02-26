package repository

// RepositoryResult is a struct to wrap repository result
// so its easy to use it in channel
type RepositoryResult struct {
	Result interface{}
	Error  error
}

// EventWrapper is used to wrap the event interface with its struct name,
// so it will be easier to unmarshal later
type EventWrapper struct {
	EventName string
	EventData interface{}
}
