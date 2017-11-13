// Package routing provides the list of functions for each HTTP routing
package routing

type RequestError struct {
	Message string
}

const ErrorFarmID = "?farm_id query parameter is not found or empty."
