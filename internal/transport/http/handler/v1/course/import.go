package course

import (
	"icms/internal/model"
	"icms/internal/transport/http/request"
	"icms/pkg/response"
	"net/http"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func (h *CourseHandler) Import(c *gin.Context) {
	req := request.CourseImportRequest{}
	if ok := request.BindAndValidate(c, &req); !ok {
		return
	}

	file, err := req.XLSX.Open()

	if err != nil {
		response.Abort500(c, err.Error())
		return
	}

	f, err := excelize.OpenReader(file)

	if err != nil {
		response.Abort500(c, err.Error())
		return
	}

	rows := f.GetRows(f.GetSheetName(1))
	colsIndex := map[string]int{}
	coursesHash := map[string]*model.Course{}

	for colIdx, col := range rows[0] {
		colsIndex[col] = colIdx
	}

	for i := 1; i < len(rows); i++ {
		row := rows[i]

		course := &model.Course{
			Code:      getIndexOfRow(&row, colsIndex["COURSE CODE"]),
			Year:      cast.ToInt(getIndexOfRow(&row, colsIndex["TERM"])[:4]),
			Section:   getIndexOfRow(&row, colsIndex["CLASS SECTION"]),
			Title:     getIndexOfRow(&row, colsIndex["COURSE TITLE"]),
			Timeslots: model.CourseTimeslots{},
		}

		identifier := course.Code + cast.ToString(course.Year) + course.Section

		if _, ok := coursesHash[identifier]; ok {
			course = coursesHash[identifier]
		} else {
			coursesHash[identifier] = course
		}

		course.Timeslots = buildSlots(course.Timeslots, &row, colsIndex)
	}

	courses := []model.Course{}
	for _, course := range coursesHash {
		courses = append(courses, *course)
	}

	err = h.courseRepo.UpsertCourse(courses)

	if err != nil {
		response.Abort500(c, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, true, "", nil)
}

func parseDate(row *[]string, index int) string {
	return getIndexOfRow(row, index)
}

func getIndexOfRow(row *[]string, index int) string {
	if len(*row) > index {
		return (*row)[index]
	} else {
		return ""
	}
}

func buildDay(row *[]string, colsIndex map[string]int) uint8 {
	if getIndexOfRow(row, colsIndex["MON"]) == "MON" {
		return 1
	} else if getIndexOfRow(row, colsIndex["TUE"]) == "TUE" {
		return 2
	} else if getIndexOfRow(row, colsIndex["WED"]) == "WED" {
		return 3
	} else if getIndexOfRow(row, colsIndex["THU"]) == "THU" {
		return 4
	} else if getIndexOfRow(row, colsIndex["FRI"]) == "FRI" {
		return 5
	} else if getIndexOfRow(row, colsIndex["SAT"]) == "SAT" {
		return 6
	} else if getIndexOfRow(row, colsIndex["SUN"]) == "SUN" {
		return 7
	}
	return 8
}

func buildSlots(slots model.CourseTimeslots, row *[]string, colsIndex map[string]int) model.CourseTimeslots {
	slot := model.CourseTimeslot{
		Day:       buildDay(row, colsIndex),
		Venue:     getIndexOfRow(row, colsIndex["VENUE"]),
		StartDate: parseDate(row, colsIndex["START DATE"]),
		EndDate:   parseDate(row, colsIndex["END DATE"]),
		StartTime: getIndexOfRow(row, colsIndex["START TIME"]),
		EndTime:   getIndexOfRow(row, colsIndex["END TIME"]),
	}

	return append(slots, slot)
}
