package repository

import (
	"github.com/google/wire"

	"icms/internal/repository/activity"
	"icms/internal/repository/course"
	"icms/internal/repository/enrolment"
	"icms/internal/repository/message"
	"icms/internal/repository/user"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(activity.New),
	wire.NewSet(course.New),
	wire.NewSet(enrolment.New),
	wire.NewSet(message.New),
	wire.NewSet(user.New),
)
