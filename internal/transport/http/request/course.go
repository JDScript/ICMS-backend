package request

type CoursePaginateRequest struct {
	PaginationRequest
	Search *string `query:"search"`
}
type CourseGetRequest struct {
	CourseID int64 `path:"courseId"`
}

type CourseSectionsGetRequest struct {
	CourseGetRequest
}

type CourseMessagesGetRequest struct {
	CourseGetRequest
}
