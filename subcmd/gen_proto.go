package subcmd

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type GenProto struct {
	domain string
}

func NewGenProto(domain string) *GenProto {
	return &GenProto{
		domain: domain,
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
		fmt.Println("---------")
		fmt.Println(v)
		fmt.Println("---------")
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

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("message %s {\n", p.Name))
	for i, v := range p.Field {
		sb.WriteString(fmt.Sprintf("    %s %s = %d;\n", getFieldType(v.Type), v.Name, i+1))
	}
	sb.WriteString("}")

	return sb.String(), nil
}

func (g *GenProto) ListStruct() ([]Struct, error) {
	result := make([]Struct, 0)

	root := fmt.Sprintf(fmt.Sprintf("internal/domain/%s", g.domain))
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
		fmt.Println(err)
	}

	var parseField func(input *ast.Field) (Field, error)

	parseField = func(field *ast.Field) (Field, error) {
		f := Field{}

		// 获取字段名称
		if len(field.Names) > 0 {
			f.Name = field.Names[0].Name
		} else {
			f.Anonymous = true
		}

		// 获取字段类型
		if fieldType, ok := field.Type.(*ast.Ident); ok {
			f.Type = FieldType{
				Name: fieldType.Name,
			}
		} else if fieldType, ok := field.Type.(*ast.SelectorExpr); ok {
			f.Type = FieldType{
				Name: fmt.Sprintf("%s.%s", fieldType.X.(*ast.Ident).Name, fieldType.Sel.Name),
			}
		} else if fieldType, ok := field.Type.(*ast.StarExpr); ok {
			if ident, ok := fieldType.X.(*ast.Ident); ok {
				f.Type = FieldType{
					Name: ident.Name,
				}
			}
		} else if fieldType, ok := field.Type.(*ast.ArrayType); ok {
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
				case *ast.StarExpr:
					elt = v.X
				case *ast.ArrayType:
					curr.Items = &FieldType{
						Name: "array",
					}
					curr = curr.Items
					elt = v.Elt
				default:
					break outer
				}
			}
		} else if fieldType, ok := field.Type.(*ast.MapType); ok {
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
				case *ast.StarExpr:
					elt = v.X
				case *ast.ArrayType:
					curr.KeyItem = &FieldType{
						Name: "array",
					}
					curr = curr.KeyItem
					elt = v.Elt
				default:
					break keyOut
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
				case *ast.StarExpr:
					elt = v.X
				case *ast.ArrayType:
					curr.ValueItem = &FieldType{
						Name: "array",
					}
					curr = curr.ValueItem
					elt = v.Elt
				default:
					break valOut
				}
			}
		} else {
			return Field{}, errors.New("field unknown")
		}

		return f, nil
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

							// 遍历结构体字段
							for _, field := range structType.Fields.List {
								f, err := parseField(field)
								if err != nil {
									return nil, err
								}

								s.Field = append(s.Field, f)
							}

							result = append(result, s)
						}
					}
				}
			}
		}
	}

	return result, nil
}
