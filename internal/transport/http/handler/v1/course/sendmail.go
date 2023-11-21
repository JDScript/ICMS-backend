package course

import (
	"icms/internal/model"
	"icms/internal/transport/http/request"
	"icms/pkg/auth"
	"icms/pkg/response"
	"io"
	"os"
	"time"

	"html/template"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

var mailDialer = gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"))
var mailTemplate = template.Must(template.ParseFiles("templates/layout.html"))

func (hander *CourseHandler) SendMail(c *gin.Context) {
	req := request.CourseSendMailRequest{}
	if ok := request.BindAndValidate(c, &req); !ok {
		return
	}
	user := auth.CurrentUser(c)
	if enrol := hander.enrolmentRepo.IsEnrolledInCourse(user.ID, req.CourseID); enrol == nil {
		response.Abort403(c, "You haven't enrolled in this course")
		return
	}
	course := *hander.courseRepo.GetCourse(req.CourseID)
	hk, _ := time.LoadLocation("Asia/Hong_Kong")
	now := time.Now().In(hk)
	weekday := now.Weekday()
	date := now.Format(time.DateOnly)
	found := false
	var upcomingTimeslot model.CourseTimeslot
	var upcomingTime time.Time
	for _, timeslot := range course.Timeslots {
		if timeslot.StartDate <= date && date <= timeslot.EndDate {
			slottime, _ := time.ParseInLocation(time.DateTime, date+" "+timeslot.StartTime+":00", hk)
			slottime = slottime.Add(time.Hour * 24 * (time.Duration(timeslot.Day%7) - time.Duration(weekday)))
			if now.Before(slottime) && (!found || slottime.Before(upcomingTime)) {
				upcomingTimeslot = timeslot
				upcomingTime = slottime
				found = true
			}
		}
	}
	c.Set("page_size", "5")
	c.Set("page", "1")
	messages := hander.messageRepo.GetCourseMessages(c, req.CourseID, user.ID).List.([]model.CourseMessage)
	sections := hander.courseRepo.GetCourseContents(req.CourseID)
	mail := gomail.NewMessage()
	mail.SetHeader("From", os.Getenv("SMTP_USER"))
	mail.SetHeader("To", user.Email)
	mail.SetHeader("Subject", "Your course information is ready")
	var info = struct {
		User     model.User
		Course   model.Course
		Timeslot model.CourseTimeslot
		Messages []model.CourseMessage
		Sections []model.CourseSection
	}{
		User:     *user,
		Course:   course,
		Timeslot: upcomingTimeslot,
		Messages: messages,
		Sections: sections,
	}
	mail.SetBody("text/plain", "No HTML Support")
	mail.AddAlternativeWriter("text/html", func(w io.Writer) error {
		return mailTemplate.Execute(w, info)
	})
	mailDialer.DialAndSend(mail)
	response.JSON(c, 200, true, "Mail successfully sent", nil)
}
