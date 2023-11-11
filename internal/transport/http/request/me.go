package request

type MeEnrolCourseRequest struct {
	CourseID int64 `json:"course_id" vd:"@:$>0; msg:'Invalid course_id'"`
}
