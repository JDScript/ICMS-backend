package enum

type UserIdentity string

const (
	UserIdentity_Student UserIdentity = "student"
	UserIdentity_Teacher UserIdentity = "teacher"
	UserIdentity_Admin   UserIdentity = "admin"
)
