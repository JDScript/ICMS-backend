package course

import (
	"bytes"
	"fmt"
	"html/template"
	"icms/internal/config"
	"icms/internal/model"
	"icms/internal/transport/http/request"
	"icms/pkg/auth"
	"icms/pkg/path"
	"icms/pkg/response"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

type CourseEmailData struct {
	Subject        string
	User           model.User
	Course         model.Course
	LatestMessages []model.CourseMessage
	LatestModules  []model.CourseModule
	Start          string
	End            string
	Venue          string
}

func (handler *CourseHandler) SendMail(c *gin.Context) {
	req := request.CourseSendMailRequest{}
	if ok := request.BindAndValidate(c, &req); !ok {
		return
	}

	user := auth.CurrentUser(c)
	if enrol := handler.enrolmentRepo.IsEnrolledInCourse(user.ID, req.CourseID); enrol == nil {
		response.Abort403(c, "You haven't enrolled in this course")
		return
	}

	course := handler.courseRepo.GetCourse(req.CourseID)
	hk, _ := time.LoadLocation("Asia/Hong_Kong")
	now := time.Now().In(hk)

	var upcomingStartTime time.Time
	var upcomingEndTime time.Time
	var venue string
	found := false

	for _, timeslot := range course.Timeslots {
		start, _ := time.ParseInLocation(
			"2006-01-02",
			timeslot.StartDate,
			hk,
		)
		end, _ := time.ParseInLocation(
			"2006-01-02",
			timeslot.EndDate,
			hk,
		)
		end = end.Add(time.Hour * 24)

		if now.After(start) && now.Before(end) {
			// 本周的 slot time
			slottime, _ := time.ParseInLocation(time.DateTime, now.Format(time.DateOnly)+" "+timeslot.StartTime+":00", hk)
			slottime = slottime.Add(time.Hour * 24 * (time.Duration(timeslot.Day%7) - time.Duration(now.Weekday())))

			if slottime.After(now) && (!found || slottime.Before(upcomingStartTime)) {
				upcomingStartTime = slottime
				upcomingEndTime, _ = time.ParseInLocation(time.DateTime, slottime.Format(time.DateOnly)+" "+timeslot.EndTime+":00", hk)
				venue = timeslot.Venue
				found = true
			}
		}
	}

	latestMessages := handler.messageRepo.GetLatestCourseMessages(req.CourseID, 3)
	latestModules := handler.courseRepo.GetLatestCourseModules(req.CourseID, 3)
	emailData := CourseEmailData{
		Subject:        fmt.Sprintf("%s [Section, %s] - ICMS Course Email", course.Code, course.Section),
		Course:         *course,
		User:           *user,
		LatestMessages: latestMessages,
		LatestModules:  latestModules,
		Start:          upcomingStartTime.Format("2006-01-02 15:04") + " HKT",
		End:            upcomingEndTime.Format("2006-01-02 15:04") + " HKT",
		Venue:          venue,
	}

	// Load Template and render
	emailTemplate := template.Must(template.ParseFiles(filepath.Join(path.RootPath(), "external/templates/email.html")))
	var renderBuffer bytes.Buffer
	emailTemplate.Execute(&renderBuffer, emailData)

	// Mail
	mail := gomail.NewMessage()
	mail.SetHeader("From", config.GlobalConfig.Email.User)
	mail.SetHeader("To", user.Email)
	mail.SetHeader("Subject", emailData.Subject)
	mail.SetBody("text/html", renderBuffer.String())

	mailDialer := gomail.NewDialer(
		config.GlobalConfig.Email.Host,
		config.GlobalConfig.Email.Port,
		config.GlobalConfig.Email.User,
		config.GlobalConfig.Email.Password,
	)

	if err := mailDialer.DialAndSend(mail); err != nil {
		response.Abort500(c, err.Error())
		return
	}

	response.JSON(c, 200, true, "Mail successfully sent", nil)
}
