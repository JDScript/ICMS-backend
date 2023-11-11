package enum

type ActivityType string

const (
	Activity_Register            ActivityType = "register"
	Activity_Login               ActivityType = "login"
	Activity_Get_My_Profile      ActivityType = "get_my_profile"
	Activity_Get_My_Activities   ActivityType = "get_my_activities"
	Activity_Get_My_Enrolments   ActivityType = "get_my_enrolments"
	Activity_Greate_My_Enrolment ActivityType = "create_my_enrolment"
	Activity_Get_My_Messages     ActivityType = "get_my_messages"
	Activity_Search_All_Courses  ActivityType = "search_all_courses"
	Activity_Get_Course_Contents ActivityType = "get_course_contents"
	Activity_Get_Course_Messages ActivityType = "get_course_messages"
	Activity_Unknown             ActivityType = "unknown"
)
