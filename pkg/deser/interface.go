package deser

type DeserFunc = func([]byte, *map[string]any) error
