package subcmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type GenProto struct {
	domain string
	name   string
}

func NewGenProto(domain string, name string) *GenProto {
	return &GenProto{
		domain: domain,
		name:   name,
	}
}

func (g *GenProto) Process() error {
	data, err := g.ListStruct()
	if err != nil {
		return err
	}

	result := make([]string, 0)

	for _, v := range data {
		msg, err := g.GenMessage(v)
		if err != nil {
			return err
		}
		result = append(result, msg)
	}

	for _, v := range result {
		fmt.Println(v)
	}

	return nil
}

func (g *GenProto) GenMessage(p Struct) (string, error) {
	var getFieldType func(typ FieldType, in ...string) string

	getFieldType = func(typ FieldType, in ...string) string {
		result := make([]string, 0)

	loop:
		switch typ.Name {
		case "int", "int8", "int16", "int32":
			result = append(result, "int32")
		case "uint", "uint8", "uint16", "uint32":
			result = append(result, "uint32")
		case "array":
			result = append(result, "repeated")
			if len(in) > 0 && in[0] == "key" {
				for typ.KeyItem != nil {
					typ = *typ.KeyItem
					goto loop
				}
			} else if len(in) > 0 && in[0] == "value" {
				for typ.ValueItem != nil {
					typ = *typ.ValueItem
					goto loop
				}
			} else {
				for typ.Items != nil {
					typ = *typ.Items
					goto loop
				}
			}
		case "map":
			result = append(result, fmt.Sprintf("map<%s, %s>", getFieldType(*typ.KeyItem, "key"), getFieldType(*typ.ValueItem, "value")))
		default:
			result = append(result, typ.Name)
		}

		return strings.Join(result, " ")
	}

	lowerFirst := func(s string) string {
		r := []rune(s)
		r[0] = unicode.ToLower(r[0])
		return string(r)
	}

	index := 6

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("message %s {\n", p.Name))

	for _, v := range p.Field {
		if v.Name == "" || g.isLower(v.Name) {
			continue
		}

		sb.WriteString(fmt.Sprintf("    %s %s = %d;", getFieldType(v.Type), lowerFirst(v.Name), index))

		if len(v.GoTag) > 0 {
			sb.WriteString(fmt.Sprintf("  //@gotags: %s", strings.Join(v.GoTag, " ")))
		}
		sb.WriteString("\n")

		index++
	}

	sb.WriteString("}")

	return sb.String(), nil
}

func (g *GenProto) ListStruct() ([]Struct, error) {
	result := make([]Struct, 0)
	root := fmt.Sprintf(fmt.Sprintf("internal/domain/%s", g.domain))
	//root := fmt.Sprintf(fmt.Sprintf("%s", g.domain))
	goFile := make([]string, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			goFile = append(goFile, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var parseField func(name string, input *ast.Field) (Field, error)

	parseField = func(name string, field *ast.Field) (Field, error) {
		f := Field{}
		if field.Tag != nil {
			//@gotags: remark:"记录id"
			goTag := make([]string, 0)
			for _, v := range strings.Fields(strings.Trim(field.Tag.Value, "`")) {
				if val := strings.Split(v, ":"); len(val) == 2 {
					if val[0] == "remark" || val[0] == "validate" {
						goTag = append(goTag, v)
					}
				}
			}
			f.GoTag = goTag
		}

		// 获取字段名称
		if len(field.Names) > 0 {
			f.Name = field.Names[0].Name
		} else {
			f.Anonymous = true
		}

		// 获取字段类型
		switch fieldType := field.Type.(type) {
		case *ast.Ident:
			f.Type = FieldType{
				Name: fieldType.Name,
			}
		case *ast.SelectorExpr:
			f.Type = FieldType{
				Name: fmt.Sprintf("%s.%s", fieldType.X.(*ast.Ident).Name, fieldType.Sel.Name),
			}
		case *ast.StarExpr:
			if ident, ok := fieldType.X.(*ast.Ident); ok {
				f.Type = FieldType{
					Name: ident.Name,
				}
			}
		case *ast.ArrayType:
			f.Type = FieldType{
				Name: "array",
			}

			curr := &f.Type

			elt := fieldType.Elt
		outer:
			for {
				switch v := elt.(type) {
				case *ast.Ident:
					curr.Items = &FieldType{
						Name: v.Name,
					}
					break outer
				case *ast.SelectorExpr:
					curr.Items = &FieldType{
						Name: fmt.Sprintf("%s.%s", v.X.(*ast.Ident).Name, v.Sel.Name),
					}
					break outer
				case *ast.StarExpr:
					elt = v.X
				case *ast.ArrayType:
					curr.Items = &FieldType{
						Name: "array",
					}
					curr = curr.Items
					elt = v.Elt
				default:
					return Field{}, fmt.Errorf("[WARN] struct[%s] field[%s] unknown not generate", name, f.Name)
				}
			}
		case *ast.MapType:
			f.Type = FieldType{
				Name: "map",
			}

			curr := &f.Type

			elt := fieldType.Key
		keyOut:
			for {
				switch v := elt.(type) {
				case *ast.Ident:
					curr.KeyItem = &FieldType{
						Name: v.Name,
					}
					break keyOut
				case *ast.SelectorExpr:
					curr.KeyItem = &FieldType{
						Name: fmt.Sprintf("%s.%s", v.X.(*ast.Ident).Name, v.Sel.Name),
					}
					break keyOut
				case *ast.StarExpr:
					elt = v.X
				case *ast.ArrayType:
					curr.KeyItem = &FieldType{
						Name: "array",
					}
					curr = curr.KeyItem
					elt = v.Elt
				default:
					return Field{}, fmt.Errorf("[WARN] struct[%s] field[%s] unknown not generate", name, f.Name)
				}
			}

			curr = &f.Type

			elt = fieldType.Value
		valOut:
			for {
				switch v := elt.(type) {
				case *ast.Ident:
					curr.ValueItem = &FieldType{
						Name: v.Name,
					}
					break valOut
				case *ast.SelectorExpr:
					curr.ValueItem = &FieldType{
						Name: fmt.Sprintf("%s.%s", v.X.(*ast.Ident).Name, v.Sel.Name),
					}
					break valOut
				case *ast.StarExpr:
					elt = v.X
				case *ast.ArrayType:
					curr.ValueItem = &FieldType{
						Name: "array",
					}
					curr = curr.ValueItem
					elt = v.Elt
				default:
					return Field{}, fmt.Errorf("[WARN] struct[%s] field[%s] unknown not generate", name, f.Name)
				}
			}
		default:
			return Field{}, fmt.Errorf("[WARN] struct[%s] field[%s] unknown not generate", name, f.Name)
		}
		return f, nil
	}

	list := make([]string, 0)
	list = append(list, g.name)
	set := make(map[string]struct{})

	for len(list) > 0 {
		process := append([]string{}, list...)
		list = make([]string, 0)

		for _, name := range process {
			if _, exist := set[name]; !exist {
				set[name] = struct{}{}
			} else {
				continue
			}

			for _, v := range goFile {
				fileSet := token.NewFileSet()

				file, err := parser.ParseFile(fileSet, v, nil, parser.ParseComments)
				if err != nil {
					return nil, err
				}

				// 遍历文件中的所有声明
				for _, decl := range file.Decls {
					// 如果声明是结构体类型
					if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
						for _, spec := range genDecl.Specs {
							if typeSpec, ok := spec.(*ast.TypeSpec); ok {
								// 获取结构体类型
								if structType, ok := typeSpec.Type.(*ast.StructType); ok {
									// 获取结构体名称
									s := Struct{
										Name: typeSpec.Name.Name,
									}

									if s.Name != name {
										continue
									}

									// 遍历结构体字段
									for _, field := range structType.Fields.List {
										f, err := parseField(s.Name, field)
										if err != nil {
											fmt.Println(err.Error())
											continue
										}
										s.Field = append(s.Field, f)

										if !g.isBasic(f.Type.Name) {
											ht := f.Type.HoldType()

											if len(ht) > 0 {
												for _, h := range ht {
													list = append(list, h.Name)
												}
											} else {
												list = append(list, f.Type.Name)
											}
										}
									}

									result = append(result, s)
								}
							}
						}
					}
				}
			}
		}
	}

	return result, nil
}

func (g *GenProto) isLower(s string) bool {
	r := []rune(s)
	return len(r) > 0 && unicode.IsLower(r[0])
}

func (g *GenProto) isBasic(s string) bool {
	switch s {
	case "int", "int8", "int16", "int32", "int64":
		return true
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return true
	case "float32", "float64":
		return true
	case "bool":
		return true
	case "string":
		return true
	default:
		return false
	}
}
