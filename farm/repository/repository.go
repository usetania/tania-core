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

func getRandomUID() string {
	rand.Seed(time.Now().UnixNano())
	uid := rand.Int()
	return strconv.Itoa(uid)
}
