package category

type ListCategoriesRequest struct {
	Page     int     `form:"page,default=1" binding:"min=1"`
	PageSize int     `form:"page_size,default=20" binding:"min=1,max=100"`
	Keyword  string  `form:"keyword" binding:"omitempty,max=64"`
	Status   *int    `form:"status" binding:"omitempty,oneof=0 1"`
	ParentID *uint64 `form:"parent_id"`
}

type CategoryTreeRequest struct {
	Status   *int    `form:"status" binding:"omitempty,oneof=0 1"`
	ParentID *uint64 `form:"parent_id"`
}

type CreateCategoryRequest struct {
	ParentID uint64 `json:"parent_id"`
	Name     string `json:"name" binding:"required,min=1,max=64"`
	Code     string `json:"code" binding:"required,min=2,max=64"`
	Sort     int    `json:"sort" binding:"omitempty,min=0,max=1000000"`
	Status   *int   `json:"status" binding:"omitempty,oneof=0 1"`
	Icon     string `json:"icon" binding:"omitempty,max=128"`
	Remark   string `json:"remark" binding:"omitempty,max=255"`
}

type UpdateCategoryRequest struct {
	ParentID uint64 `json:"parent_id"`
	Name     string `json:"name" binding:"required,min=1,max=64"`
	Code     string `json:"code" binding:"required,min=2,max=64"`
	Sort     int    `json:"sort" binding:"omitempty,min=0,max=1000000"`
	Status   int    `json:"status" binding:"required,oneof=0 1"`
	Icon     string `json:"icon" binding:"omitempty,max=128"`
	Remark   string `json:"remark" binding:"omitempty,max=255"`
}

type CategoryItem struct {
	ID       uint64 `json:"id"`
	ParentID uint64 `json:"parent_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status"`
	Icon     string `json:"icon"`
	Remark   string `json:"remark"`
}

type ListCategoriesResult struct {
	List     []CategoryItem `json:"list"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

type CategoryTreeItem struct {
	ID       uint64             `json:"id"`
	ParentID uint64             `json:"parent_id"`
	Name     string             `json:"name"`
	Code     string             `json:"code"`
	Sort     int                `json:"sort"`
	Status   int                `json:"status"`
	Icon     string             `json:"icon"`
	Remark   string             `json:"remark"`
	Children []CategoryTreeItem `json:"children"`
}
