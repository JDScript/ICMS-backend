package request

type AuthenticationLoginRequest struct {
	Face string `json:"face" vd:"@:regexp('^data:image\\/jpeg;base64,[A-Za-z0-9+/]+={0,2}$'); msg:'Face image should be encoded using jpeg format'"`
}
