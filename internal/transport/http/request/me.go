package request

type MeEnrolCourseRequest struct {
	CourseID int64 `json:"course_id" vd:"@:$>0; msg:'Invalid course_id'"`
}

type MeGetMessagesRequest struct {
	PaginationRequest
	CourseId *int64 `query:"course_id"`
	Unread   bool   `query:"unread" default:"false"`
	Order    string `query:"order" vd:"@:mblen($)==0; msg:'No order should be provided'" default:""`
	Sort     string `query:"sort" vd:"@:mblen($)==0; msg:'No sort should be provided'" default:""`
}

type MeReadMessagesRequest struct {
	MessagesID []int64 `json:"messages_id" default:"[]"`
}
