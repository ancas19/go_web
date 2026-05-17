package meta

import (
	"os"
	"strconv"
)

type (
	Meta struct {
		TotalCount int64 `json:"totalCount"`
		Page       int64 `json:"page"`
		PerPage    int64 `json:"perPage"`
	}
)

func New(total, page, perPage int64) (*Meta, error) {

	if perPage <= 0 {
		maxValue, err := strconv.ParseInt(os.Getenv("PAGINATOR_LIMIT_DEFAULT"), 10, 64)
		if err != nil {
			return nil, err
		}
		perPage = maxValue
	}
	var pageCount int64 = 0
	if total >= 0 {
		pageCount = (total + perPage - 1) / perPage
		if page > pageCount {
			page = pageCount
		}
	}
	if page < 1 {
		page = 1
	}
	return &Meta{
		TotalCount: total,
		Page:       page,
		PerPage:    perPage,
	}, nil
}

func (m *Meta) Offset() int64 {
	return (m.Page - 1) * m.PerPage
}

func (m *Meta) Limit() int64 {
	return m.PerPage
}
