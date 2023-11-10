package authentication

import (
	"encoding/base64"
	"icms/internal/component/jwt"
	"icms/internal/transport/http/request"
	"icms/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func (handler *AuthenticationHandler) Login(c *gin.Context) {
	req := request.AuthenticationLoginRequest{}
	if ok := request.BindAndValidate(c, &req); !ok {
		return
	}

	// Encode JPEG image
	base64Img := strings.TrimPrefix(req.Face, "data:image/jpeg;base64,")
	imgData, err := base64.StdEncoding.DecodeString(base64Img)
	if err != nil {
		response.BadRequest(c, err, "Error when decoding facial image")
		return
	}

	face, err := handler.facial.RecognizeCNN(imgData)
	if err != nil {
		response.Abort500(c, err.Error())
		return
	}

	if face == nil {
		response.AbortJSON(c, http.StatusBadRequest, false, "No face detected", nil)
		return
	}

	userId := handler.facial.Classify(face.Descriptor)

	user := handler.userRepo.GetByID(cast.ToString(userId))
	if user.ID == 0 {
		response.AbortJSON(c, http.StatusBadRequest, false, "No matched user", nil)
		return
	}

	token, err := handler.jwt.IssueToken(&jwt.JwtCustomClaims{
		UserId: cast.ToString(userId),
	})

	if err != nil {
		response.Abort500(c, err.Error())
		return
	}

	// For activity logging middleware
	c.Set("current_user", user)
	c.Set("current_user_id", cast.ToString(user.ID))

	response.JSON(c, http.StatusCreated, true, "Login successfully", token)
}
