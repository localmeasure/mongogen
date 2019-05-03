// Copyright 2019 Local Measure. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mongogen

import (
	"log"
	"strings"
)

var (
	bsonMap = map[string]string{
		"id":     "primitive.ObjectID",
		"[]id":   "[]primitive.ObjectID",
		"time":   "time.Time",
		"[]time": "[]time.Time",
	}
	typOps = map[string][]string{
		"primitive.ObjectID": []string{"eq", "ne", "in", "nin", "gt", "gte", "lt", "lte"},
		"string":             []string{"eq", "ne", "in", "nin"},
		"bool":               []string{"eq", "ne"},
		"int":                []string{"eq", "ne", "in", "nin", "gt", "gte", "lt", "lte"},
		"float64":            []string{"eq", "ne", "in", "nin", "gt", "gte", "lt", "lte"},
		"time.Time":          []string{"gt", "gte", "lt", "lte"},

		// multikey indexes
		"[]primitive.ObjectID": []string{"eq", "ne", "in", "nin", "gt", "gte", "lt", "lte", "all"},
		"[]string":             []string{"eq", "ne", "in", "nin", "all"},
		"[]int":                []string{"eq", "ne", "in", "nin", "gt", "gte", "lt", "lte", "all"},
		"[]float64":            []string{"eq", "ne", "in", "nin", "gt", "gte", "lt", "lte", "all"},
		"[]time.Time":          []string{"gt", "gte", "lt", "lte", "all"},
	}
	pkgImports = map[string]string{
		"primitive.ObjectID":   "go.mongodb.org/mongo-driver/bson/primitive",
		"[]primitive.ObjectID": "go.mongodb.org/mongo-driver/bson/primitive",
		"time.Time":            "time",
		"[]time.Time":          "time",
	}
	goKeywords = map[string]struct{}{
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
)

type index struct {
	name string
	keys []indexKey
}

type indexKey struct {
	name   string
	goname string
	typ    string
}

func analyze(parsed []string, prefix string) pkg {
	pkg := pkg{
		imports: map[string]struct{}{
			"go.mongodb.org/mongo-driver/bson":           struct{}{},
			"go.mongodb.org/mongo-driver/bson/primitive": struct{}{},
			"go.mongodb.org/mongo-driver/mongo":          struct{}{},
			"go.mongodb.org/mongo-driver/mongo/options":  struct{}{},
			"context": struct{}{},
		},
	}
	var indexes []index
	for i := 0; i < len(parsed); i++ {
		var idx index
		parsedKeys := strings.Split(parsed[i], "+")
		first := true
		for n := 0; n < len(parsedKeys); n++ {
			spec := strings.Split(parsedKeys[n], ":")
			if len(spec) < 2 {
				first = false
				continue
			}
			typ, ok := bsonMap[spec[1]]
			if !ok {
				typ = spec[1]
			}
			_, ok = typOps[typ]
			if first && !ok {
				log.Fatal("Fatal: type of first key is not supported")
			}
			_, ok = pkgImports[typ]
			if ok {
				// import necessary packages
				pkg.imports[pkgImports[typ]] = struct{}{}
			}
			goname := toCamelCase(spec[0], false)
			idx.name += toCamelCase(goname, true)
			idx.keys = append(idx.keys, indexKey{
				name:   spec[0],
				goname: escapeGoKeyword(goname),
				typ:    typ,
			})
			first = false
		}
		idx.name = "use" + idx.name
		indexes = append(indexes, idx)
	}
	pkg.indexes = indexes
	return pkg
}

func escapeGoKeyword(key string) string {
	if _, matched := goKeywords[key]; matched {
		return key + "Arg"
	}
	return key
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
