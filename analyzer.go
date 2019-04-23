package mongogen

import (
	"log"
	"strconv"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/jinzhu/inflection"
	"go.mongodb.org/mongo-driver/mongo"
)

type Analyzer struct {
	sess   *mgo.Session
	db     *mongo.Database
	dbName string
	checkc chan typCheck
	donec  chan typCheck
}

func NewAnalyzer(sess *mgo.Session, db *mongo.Database, dbName string) *Analyzer {
	return &Analyzer{sess: sess, db: db, dbName: dbName}
}

type Pkg struct {
	Imports  []string `json:"imports"`
	Types    []Typ    `json:"services"`
	Database string   `json:"database"`
	Consts   [][2]string
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
	var pkg Pkg
	colNames, err := a.DB().CollectionNames()
	if err != nil {
		return pkg, err
	}
	for _, colName := range colNames {
		log.Printf("analyze collection %s\n", colName)
		indexes, err := a.DB().C(colName).Indexes()
		if err != nil {
			return pkg, err
		}
		indexes = rmMGOPrefix(indexes)
		camelColName := toCamelCase(colName, true)
		service := Typ{
			Collection: colName,
			Name:       camelColName,
			Plural:     inflection.Plural(camelColName),
			Singular:   inflection.Singular(camelColName),
		}
		var methods []Method
		methodSet := make(map[string]struct{})
		for nth, idx := range indexes {
			log.Printf("\tanalyze index %s\n", idx.Name)
			constant := service.Plural + "Idx" + strconv.FormatInt(int64(nth+1), 10)
			pkg.Consts = append(pkg.Consts, [2]string{constant, idx.Name})
			if len(idx.Key) == 1 && idx.Key[0] == "_id" {
				continue
			}
			a.typesOf(colName, idx)
			mName := service.Singular + "With"
			mArgs := []Argument{}
			for i := 0; i < len(idx.Key); i++ {
				mName += toCamelCase(idx.Key[i], true)
				typReg.RLock()
				argType, ok := typReg.ref[colName+idx.Key[i]]
				typReg.RUnlock()
				if !ok {
					argType = "interface{}"
				}
				mArgs = append(mArgs, Argument{
					QueryName: idx.Key[i],
					ArgName:   escapeGoKeyword(toCamelCase(idx.Key[i], false)),
					ArgType:   argType,
				})
				if _, ok := methodSet[mName]; !ok {
					methods = append(methods, Method{
						Name: mName,
						Args: mArgs,
						Hint: constant,
					})
					methodSet[mName] = struct{}{}
				}
			}
		}
		service.Methods = methods
		pkg.Types = append(pkg.Types, service)
	}

	if len(pkg.Types) > 0 {
		// for now
		pkg.Database = a.dbName
		pkg.Imports = []string{
			"context",
			"time",
			"go.mongodb.org/mongo-driver/bson",
			"go.mongodb.org/mongo-driver/bson/primitive",
			"go.mongodb.org/mongo-driver/mongo",
			"go.mongodb.org/mongo-driver/mongo/options",
		}
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

func rmMGOPrefix(indexes []mgo.Index) []mgo.Index {
	cp := make([]mgo.Index, len(indexes))
	for n, idx := range indexes {
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
		}
		cp[n] = idx
	}
	return cp
}
