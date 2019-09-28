package reflectutil

type ctx struct {
	*Config
	encoder ValEncoder
}

type StructField struct {
	Encoder *StructFieldEncoder
	Name    string
}
