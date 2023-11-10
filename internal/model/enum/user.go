package enum

type UserIdentity string

const (
	UserIdentity_Student EnrolmentIdentity = "student"
	UserIdentity_Teacher EnrolmentIdentity = "teacher"
	UserIdentity_Admin   EnrolmentIdentity = "admin"
)
