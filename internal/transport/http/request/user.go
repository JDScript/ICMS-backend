package request

type UserCreateRequest struct {
	Name        string         `json:"name" vd:"@:mblen($)>0&&mblen($)<=64; msg:'Name should be between 1~64 characters'"`
	Email       string         `json:"email" vd:"@:email($); msg:'Invalid email address'"`
	Descriptors [][128]float32 `json:"descriptors" vd:"@:len($)>=50; msg:'Need at least 50 descriptors'"`
	// Code        string         `json:"code" vd:"@:mblen($)==10; msg:'Code should be exactly 10 characters'"`
}
