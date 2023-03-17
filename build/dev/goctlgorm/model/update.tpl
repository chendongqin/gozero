
func (ins *{{.upperStartCamelObject}}) GetPkValue() interface{} {
    return ins.{{.upperStartCamelPrimaryKey}}
}

func (ins *{{.upperStartCamelObject}}) SetTime() {
    hasPkValue := false
	val := ins.GetPkValue()
	switch val.(type) {
	case int:
		hasPkValue = val.(int) > 0
	case int64:
		hasPkValue = val.(int64) > 0
	case uint8:
		hasPkValue = val.(uint8) > 0
	case int32:
		hasPkValue = val.(int32) > 0
	case string:
		hasPkValue = val.(string) != ""
	}
	if !hasPkValue {
    	ins.CreatedAt = time.Now()
    }
    ins.UpdatedAt = time.Now()

    return
}