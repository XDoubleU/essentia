package main

type Data struct {
	name string
}

type DataRepository struct {
	data map[string]Data
}

func (r DataRepository) GetPaged(pageIndex int, pageSize int) []Data {
	page := []Data{}

	keys := make([]string, 0, len(r.data))
	for k := range r.data {
		keys = append(keys, k)
	}

	for i := pageIndex * pageSize; i < len(r.data) && len(page) < pageSize; i++ {
		page = append(page, r.data[keys[i]])
	}

	return page
}

func (r DataRepository) GetSingle(id string) *Data {
	v, ok := r.data[id]
	if !ok {
		return nil
	}

	return &v
}

func (r DataRepository) Create() *Data {
	return nil
}

func (r DataRepository) Update() *Data {
	return nil
}

func (r DataRepository) Delete(id string) *Data {
	v, ok := r.data[id]
	if !ok {
		return nil
	}

	delete(r.data, id)

	return &v
}
