package user

import (
	"icms/internal/model"
	"icms/internal/transport/http/request"
	"icms/pkg/response"
	"net/http"

	"github.com/Kagami/go-face"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func (handler *UserHandler) Create(c *gin.Context) {
	req := request.UserCreateRequest{}
	if ok := request.BindAndValidate(c, &req); !ok {
		return
	}

	rawDescriptors := make([]face.Descriptor, len(req.Descriptors))
	descriptors := make([]model.FacialDescriptor, len(req.Descriptors))
	for idx := 0; idx < len(req.Descriptors); idx++ {
		descriptor := face.Descriptor(req.Descriptors[idx])

		userId := handler.facial.Classify(descriptor)
		if userId >= 0 {
			response.AbortJSON(c, http.StatusBadRequest, false, "Find user with similar face, please try again", nil)
			return
		}

		rawDescriptors[idx] = descriptor
		descriptors[idx] = model.FacialDescriptor{
			FacialDescriptor: req.Descriptors[idx],
		}
	}

	user := model.User{
		Name:        req.Name,
		Email:       req.Email,
		Descriptors: descriptors,
	}

	err := handler.userRepo.Create(&user)
	if err != nil {
		response.Abort500(c, err.Error())
		return
	}

	userIds := make([]int32, len(rawDescriptors))
	for idx := 0; idx < len(rawDescriptors); idx++ {
		userIds[idx] = user.ID
	}

	handler.facial.AddSample(rawDescriptors, userIds)

	// For activity logging middleware
	c.Set("current_user", &user)
	c.Set("current_user_id", cast.ToString(user.ID))

	response.JSON(c, http.StatusOK, true, "User created successfully", user)
}
