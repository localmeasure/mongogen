// Copyright 2019 Local Measure. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mongogen

import (
	"flag"
	"log"
	"path/filepath"
	"strings"

	"github.com/jinzhu/inflection"
	"golang.org/x/tools/go/packages"
)

var (
	bsonMap = map[string]string{
		"id":     "primitive.ObjectID",
		"[]id":   "[]primitive.ObjectID",
		"time":   "time.Time",
		"[]time": "[]time.Time",
	}
	pkgImports = map[string]string{
		"id":     "go.mongodb.org/mongo-driver/bson/primitive",
		"[]id":   "go.mongodb.org/mongo-driver/bson/primitive",
		"time":   "time",
		"[]time": "time",
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

func analyze(indexes []string, prefix string) pkg {
	pkg := pkg{
		name: getPkgName(),
		imports: map[string]struct{}{
			"context":                                   struct{}{},
			"go.mongodb.org/mongo-driver/bson":          struct{}{},
			"go.mongodb.org/mongo-driver/mongo":         struct{}{},
			"go.mongodb.org/mongo-driver/mongo/options": struct{}{},
		},
	}
	methods := []method{}
	for i := 0; i < len(indexes); i++ {
		var args []methodArg
		kplus := strings.Split(indexes[i], "+")
		for n := 0; n < len(kplus); n++ {
			kcolon := strings.Split(kplus[n], ":")
			if len(kcolon) < 2 {
				continue
			}
			typ := kcolon[1]
			if _, ok := bsonMap[kcolon[1]]; ok {
				typ = bsonMap[kcolon[1]]
				pkg.imports[pkgImports[kcolon[1]]] = struct{}{}
			}
			multiple := false
			name := escapeGoKeyword(toCamelCase(kcolon[0], false))
			if typ[:2] == "[]" {
				multiple = true
				name = inflection.Plural(name)
			}
			args = append(args, methodArg{
				query:    kcolon[0],
				name:     name,
				typ:      typ,
				multiple: multiple,
			})
		}
		methodName := prefix + "With"
		var methodArgs []methodArg
		for _, arg := range args {
			methodName += toCamelCase(arg.name, true)
			methodArgs = append(methodArgs, arg)
			methods = append(methods, method{
				name: methodName,
				args: methodArgs,
			})
		}
	}
	pkg.methods = methods
	return pkg
}

func getPkgName() string {
	args := flag.Args()
	dir := "."
	if len(args) > 0 {
		// Default: process whole package in current directory.
		dir = filepath.Dir(args[0])
	}
	cfg := &packages.Config{Mode: packages.LoadSyntax, Tests: false}
	pkgs, err := packages.Load(cfg, dir)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) < 1 {
		log.Fatal("no packages found")
	}
	return pkgs[0].Name
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
