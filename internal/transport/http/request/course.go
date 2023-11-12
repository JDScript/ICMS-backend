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
	PaginationRequest
	CourseGetRequest
	Order string `query:"order" vd:"@:mblen($)==0; msg:'No order should be provided'" default:""`
	Sort  string `query:"sort" vd:"@:mblen($)==0; msg:'No sort should be provided'" default:""`
}
