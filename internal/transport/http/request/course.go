package request

type CoursePaginateRequest struct {
	PaginationRequest
	Search *string `query:"search"`
}

type CourseSectionsGetRequest struct {
	CourseID int64 `path:"courseId"`
}

type CourseMessagesGetRequest struct {
	CourseID int64 `path:"courseId"`
}
