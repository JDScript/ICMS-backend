// Package paginator 处理分页逻辑
package paginator

import (
	"math"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type Paging struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

// Paginator 分页操作类
type Paginator struct {
	PageSize  int   // 每页条数
	Current   int   // 当前页
	Offset    int   // 偏移量
	TotalPage int   // 总页数
	Total     int64 // 总条数

	Sort  string // 排序规则
	Order string // 排序顺序

	ctx   *gin.Context
	query *gorm.DB // db query 句柄
}

// 手动分页的时候这么干
func GetPaginator(c *gin.Context, total int64) *Paginator {
	p := &Paginator{
		ctx:   c,
		Total: total,
	}
	p.PageSize = p.getPageSize()
	p.TotalPage = p.getTotalPage()
	p.Current = p.getCurrent()
	p.Order = p.ctx.DefaultQuery("order", "desc")
	p.Sort = p.ctx.DefaultQuery("sort", "created_at")
	p.Offset = (p.Current - 1) * p.PageSize
	return p
}

// 自动分页的时候这么干
func Paginate(c *gin.Context, db *gorm.DB, data interface{}, stableSort ...string) Paging {
	p := &Paginator{
		query: db,
		ctx:   c,
	}
	p.PageSize = p.getPageSize()
	p.Total = p.getTotal()
	p.TotalPage = p.getTotalPage()
	p.Current = p.getCurrent()
	p.Order = p.ctx.DefaultQuery("order", "desc")
	p.Sort = p.ctx.DefaultQuery("sort", "created_at")
	p.Offset = (p.Current - 1) * p.PageSize

	if len(p.Sort) > 0 {
		p.query = p.query.Order(p.Sort + " " + p.Order)
	}

	if len(stableSort) > 0 {
		p.query = p.query.Order(stableSort)
	}

	err := p.query. // 读取关联
			Limit(p.PageSize).
			Offset(p.Offset).
			Find(data).
			Error

	if err != nil {
		return Paging{}
	}

	return Paging{
		List:  data,
		Total: p.Total,
	}
}

func (p Paginator) getPageSize() int {
	queryPerSize := p.ctx.DefaultQuery("page_size", "20")
	return cast.ToInt(queryPerSize)
}

func (p Paginator) getCurrent() int {
	// 优先取用户请求的 page
	page := cast.ToInt(p.ctx.Query("page"))
	if page <= 0 {
		// 默认为 1
		page = 1
	}
	if page > p.TotalPage {
		return p.TotalPage
	}
	return page
}

func (p *Paginator) getTotal() int64 {
	var count int64
	if err := p.query.Count(&count).Error; err != nil {
		return 0
	}
	return count
}

func (p Paginator) getTotalPage() int {
	nums := int64(math.Ceil(float64(p.Total) / float64(p.PageSize)))
	if nums == 0 {
		nums = 1
	}
	return int(nums)
}
