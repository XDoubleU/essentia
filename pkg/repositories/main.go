package repositories

type Repository[TData any, TId any] interface {
	GetPaged(pageIndex int, pageSize int) []TData
	GetSingle(id TId) *TData
	Create() *TData
	Update() *TData
	Delete(id TId) *TData
}
