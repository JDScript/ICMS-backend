package request

type ActivityPaginateRequest struct {
	PaginationRequest
	Type   *[]string `query:"type[]"`
	Method *[]string `query:"method[]"`
}
