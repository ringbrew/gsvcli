package subcmd

type Struct struct {
	Name  string
	Field []Field
	Root  bool
}

type Field struct {
	Name      string
	Type      FieldType
	Remark    string
	Required  bool
	Anonymous bool
}

type FieldType struct {
	Name      string
	Items     *FieldType
	KeyItem   *FieldType
	ValueItem *FieldType
}
