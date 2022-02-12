package requests

type GetBlocksListParam struct {
	Limit int `form:"limit" binding:"required,gte=1"`
}
