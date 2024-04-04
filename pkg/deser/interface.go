package deser

type DeserFunc = func([]byte) (interface{}, string, error)
