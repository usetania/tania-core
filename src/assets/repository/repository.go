package repository

import (
	"math/rand"
	"strconv"
	"time"
)

// RepositoryResult is a struct to wrap repository result
// so its easy to use it in channel
type RepositoryResult struct {
	Result interface{}
	Error  error
}

// GetRandomUID generates a random UID.
// Please use it before you save a struct.
func GetRandomUID() string {
	rand.Seed(time.Now().UnixNano())
	uid := rand.Int()
	return strconv.Itoa(uid)
}
