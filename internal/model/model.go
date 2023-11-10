package model

import (
	"icms/pkg/paginator"
	"icms/pkg/response"
)

type CommonTimestampField struct {
	CreatedAt int64 `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64 `gorm:"autoUpdateTime" json:"updated_at"`
}

type Response struct {
	response.Response
}

type Pagination struct {
	paginator.Paging
}
