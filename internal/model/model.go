package model

import (
	"icms/pkg/paginator"
	"icms/pkg/response"
)

type CommonTimestampField struct {
	CreatedAt int64 `gorm:"autoCreateTime:milli" json:"created_at,omitempty"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli" json:"updated_at,omitempty"`
}

type Response struct {
	response.Response
}

type Pagination struct {
	paginator.Paging
}
