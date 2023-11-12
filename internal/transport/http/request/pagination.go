package request

// 通用的分页请求，可以进行拓展/覆盖
type PaginationRequest struct {
	Page     uint64 `query:"page" vd:"@:$>0; msg:'Page number should be over 0'" default:"1"`
	PageSize uint64 `query:"page_size" vd:"@:($>0 && $<=100); msg:'Page size should be between 1 to 100'" default:"20"`
	Order    string `query:"order" vd:"@:in($,'asc','desc'); msg:'Order can be asc or desc'" default:"desc"`
	Sort     string `query:"sort" default:"created_at"`
}
