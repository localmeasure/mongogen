package mongogen

import (
	"fmt"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/jinzhu/inflection"
)

type Analyzer struct {
	sess   *mgo.Session
	dbName string
}

func NewAnalyzer(sess *mgo.Session, dbName string) *Analyzer {
	return &Analyzer{sess: sess, dbName: dbName}
}

type Pkg struct {
	Imports  []string `json:"imports"`
	Types    []Typ    `json:"services"`
	Database string   `json:"database"`
}

type Typ struct {
	Collection string   `json:"collection"`
	Name       string   `json:"name"`
	Plural     string   `json:"plural"`
	Singular   string   `json:"singular"`
	Methods    []Method `json:"methods"`
}

type Method struct {
	Name string     `json:"name"`
	Args []Argument `json:"args"`
	Hint string     `json:"hint"`
}

type Argument struct {
	QueryName string `json:"queryName"`
	ArgName   string `json:"argName"`
	ArgType   string `json:"argType"`
}

func (a *Analyzer) DB() *mgo.Database {
	return a.sess.Copy().DB(a.dbName)
}

func (a *Analyzer) Analyze() (Pkg, error) {
	colNames, err := a.DB().CollectionNames()
	if err != nil {
		return Pkg{}, err
	}

	pkg := Pkg{}
	for _, colName := range colNames {
		indexes, err := a.DB().C(colName).Indexes()
		if err != nil {
			return pkg, err
		}
		camelColName := toCamelCase(colName, true)
		service := Typ{
			Collection: colName,
			Name:       camelColName,
			Plural:     inflection.Plural(camelColName),
			Singular:   inflection.Singular(camelColName),
		}
		var methods []Method
		methodSet := make(map[string]struct{})
		typSet := make(map[string]Argument)
		for _, idx := range indexes {
			var method Method
			if len(idx.Key) == 1 && idx.Key[0] == "_id" {
				method = Method{
					Hint: idx.Name,
					Name: service.Singular + "WithID",
					Args: []Argument{{
						QueryName: "_id",
						ArgName:   "id",
						ArgType:   "primitive.ObjectID",
					}},
				}
				if _, ok := methodSet[method.Name]; !ok {
					methods = append(methods, method)
					methodSet[method.Name] = struct{}{}
				}
			} else {
				var query bson.D
				for i := 0; i < len(idx.Key); i++ {
					if strings.Index(idx.Key[i], "$text:") == 0 {
						idx.Key[i] = idx.Key[i][6:]
					}
					if strings.Index(idx.Key[i], "$2d:") == 0 {
						idx.Key[i] = idx.Key[i][4:]
					}
					if idx.Key[i][0] == '-' {
						idx.Key[i] = idx.Key[i][1:]
					}
					query = append(query, bson.DocElem{Name: idx.Key[i], Value: bson.M{"$exists": true}})
				}
				if len(query) > 0 {
					var unknown bson.M
					a.DB().C(service.Collection).Find(query).Hint(idx.Name).One(&unknown)
					for i := 0; i < len(idx.Key); i++ {
						argType := "interface{}"
						switch t := unknown[idx.Key[i]].(type) {
						case nil:
							argType = "interface{}"
						case bson.ObjectId:
							argType = "primitive.ObjectID"
						case time.Time:
							argType = "time.Time"
						case []interface{}:
							argType = "[]interface{}"
						default:
							argType = fmt.Sprintf("%T", t)
						}
						typSet[idx.Key[i]] = Argument{
							QueryName: idx.Key[i],
							ArgName:   escapeGoKeyword(toCamelCase(idx.Key[i], false)),
							ArgType:   argType,
						}
					}
				}
				for i := 0; i < len(idx.Key); i++ {
					method = Method{
						Hint: idx.Name,
						Name: service.Singular + "With",
					}
					arg := typSet[idx.Key[i]]
					method.Name += toCamelCase(idx.Key[i], true)
					method.Args = append(method.Args, arg)
					if _, ok := methodSet[method.Name]; !ok {
						methods = append(methods, method)
						methodSet[method.Name] = struct{}{}
					}
				}
			}
		}
		service.Methods = methods
		pkg.Types = append(pkg.Types, service)
	}

	if len(pkg.Types) > 0 {
		// for now
		pkg.Imports = append([]string{
			"context",
			"time",
			"go.mongodb.org/mongo-driver/bson",
			"go.mongodb.org/mongo-driver/bson/primitive",
			"go.mongodb.org/mongo-driver/mongo",
			"go.mongodb.org/mongo-driver/mongo/options",
		}, pkg.Imports...)
		pkg.Database = a.dbName
	}

	return pkg, nil
}

func toCamelCase(str string, cap bool) string {
	in := []byte(str)
	out := make([]byte, len(in))
	pos := 0
	for _, c := range in {
		if c == '_' || c == '-' || c == '.' {
			cap = true
			continue
		}
		out[pos] = c
		if c >= 'a' && c <= 'z' && cap {
			out[pos] = c - 32
		}
		cap = false
		pos++
	}
	return string(out[:pos])
}

var goKeywords = map[string]struct{}{
	"break":       struct{}{},
	"default":     struct{}{},
	"func":        struct{}{},
	"interface":   struct{}{},
	"select":      struct{}{},
	"case":        struct{}{},
	"defer":       struct{}{},
	"go":          struct{}{},
	"map":         struct{}{},
	"struct":      struct{}{},
	"chan":        struct{}{},
	"else":        struct{}{},
	"goto":        struct{}{},
	"package":     struct{}{},
	"switch":      struct{}{},
	"const":       struct{}{},
	"fallthrough": struct{}{},
	"if":          struct{}{},
	"range":       struct{}{},
	"type":        struct{}{},
	"continue":    struct{}{},
	"for":         struct{}{},
	"import":      struct{}{},
	"return":      struct{}{},
	"var":         struct{}{},
}

func escapeGoKeyword(key string) string {
	if _, matched := goKeywords[key]; matched {
		return key + "Arg"
	}
	return key
}
