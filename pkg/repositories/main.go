package repositories

type Repository interface {
	GetPaged(pageIndex int, pageSize int) []any
	GetSingle(id any) any
	Create() any
	Update() any
	Delete(id any) any
}
