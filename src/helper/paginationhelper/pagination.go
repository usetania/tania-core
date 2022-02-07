package paginationhelper

import (
	"strconv"
)

const (
	DefaultPage  int = 1
	DefaultLimit int = 10
)

// CalculatePageToOffset calculates offset based on page for query
/*
  page = 1
  limit = 10
  offset = 0

  page = 2
  limit = 10
  offset = (2-1) * 10 = 10

  page = 2
  limit = 2
  offset = (2- 1) * 2 = 2

  page = 2
  limit = 3
  offset (2 - 1) * 3 = 3
*/
func CalculatePageToOffset(page, limit int) int {
	if page == 1 {
		return 0
	}

	return (page - 1) * limit
}

func ParsePagination(page, limit string) (pageInt, limitInt int, err error) {
	pageInt = DefaultPage
	limitInt = DefaultLimit

	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil {
			return DefaultPage, DefaultLimit, err
		}
	}

	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return DefaultPage, DefaultLimit, err
		}
	}

	return pageInt, limitInt, nil
}
