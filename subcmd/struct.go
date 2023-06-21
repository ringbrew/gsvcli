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
	GoTag     []string
}

type FieldType struct {
	Name      string
	Items     *FieldType
	KeyItem   *FieldType
	ValueItem *FieldType
}

func (ft *FieldType) HoldType() []FieldType {
	result := make([]FieldType, 0)

	shadow := *ft
	curr := &shadow

	for curr.Items != nil {
		curr = curr.Items
	}
	result = append(result, *curr)

	for curr.KeyItem != nil {
		curr = curr.KeyItem
	}
	result = append(result, *curr)

	for curr.ValueItem != nil {
		curr = curr.ValueItem
	}

	result = append(result, *curr)

	return result
}
