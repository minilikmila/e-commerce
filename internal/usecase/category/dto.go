package category

type CreateCategory struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateCategoryInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type ListCategoryInput struct {
	Search   string
	Page     int
	PageSize int
}
