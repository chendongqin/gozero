package gormpool

type PageList struct {
	Page    int64
	Limit   int64
	OrderBy string
	Columns []string
	Where   string
	Values  []interface{}
}

type Conditions struct {
	Where   string
	OrderBy string
	Columns []string
	Values  []interface{}
}
