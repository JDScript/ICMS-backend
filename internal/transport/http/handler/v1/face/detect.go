package face

import (
	"encoding/base64"
	"icms/internal/transport/http/request"
	"icms/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (handler *FaceHandler) Detect(c *gin.Context) {
	req := request.FaceDetectRequest{}
	if ok := request.BindAndValidate(c, &req); !ok {
		return
	}

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

	response.JSON(c, http.StatusOK, true, "Face recognized", face)
}
