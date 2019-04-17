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
	//var importsSet = make(map[string]struct{})
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
		typSet := make(map[string]string)
		for _, idx := range indexes {
			if len(idx.Key) == 1 && idx.Key[0] == "_id" {
				continue
			}
			for i := 0; i < len(idx.Key); i++ {
				method := Method{
					Name: service.Singular + "With",
				}
				for lvl := 0; lvl <= i; lvl++ {
					queryName := idx.Key[lvl]
					if strings.Index(queryName, "$text:") == 0 {
						queryName = queryName[6:]
					}
					if strings.Index(queryName, "$2d:") == 0 {
						queryName = queryName[4:]
					}
					if queryName[0] == '-' {
						queryName = queryName[1:]
					}
					arg := Argument{
						QueryName: queryName,
						ArgName:   toCamelCase(queryName, false),
					}
					if _, ok := typSet[arg.ArgName]; !ok {
						var unknown bson.M
						a.DB().C(service.Collection).Find(bson.M{arg.QueryName: bson.M{"$exists": true}}).One(&unknown)

						switch t := unknown[arg.QueryName].(type) {
						case nil:
							arg.ArgType = "interface{}"
						case bson.ObjectId:
							arg.ArgType = "primitive.ObjectID"
						case time.Time:
							arg.ArgType = "time.Time"
						case []interface{}:
							// reformat
							arg.ArgType = "[]interface{}"
						default:
							arg.ArgType = fmt.Sprintf("%T", t)
						}
						typSet[arg.ArgName] = arg.ArgType
					} else {
						arg.ArgType = typSet[arg.ArgName]
					}
					method.Name += toCamelCase(arg.ArgName, true)
					method.Args = append(method.Args, arg)
				}

				if _, ok := methodSet[method.Name]; !ok {
					methods = append(methods, method)
					methodSet[method.Name] = struct{}{}
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
