package hzmgo

type ListObj struct {
	Total int64 `json:"total"`
	List interface{} `json:"list"`
}

func (s ListObj) GetTotal() int64 {
	return s.Total
}

func (s ListObj) GetList() interface{} {
	return s.List
}