// Copyright 2019 Local Measure. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mongogen

import (
	"bytes"
	"fmt"

	"github.com/jinzhu/inflection"
)

type Generator struct {
	buf    bytes.Buffer
	indent string
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) p(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, g.indent+format+"\n", args...)
}

func (g *Generator) in() {
	g.indent += "\t"
}

func (g *Generator) out() {
	if len(g.indent) > 0 {
		g.indent = g.indent[0 : len(g.indent)-1]
	}
}

func (g *Generator) Output() []byte {
	return g.buf.Bytes()
}

type pkg struct {
	name    string
	imports map[string]struct{}
	methods []method
}

type method struct {
	name string
	args []methodArg
}

type methodArg struct {
	query    string
	name     string
	typ      string
	multiple bool
}

func (g *Generator) Gen(collection string, indexes []string) {
	public := toCamelCase(collection, true)
	singular := inflection.Singular(public)
	plural := inflection.Plural(public)
	pkg := analyze(indexes, singular)
	g.p("// Code generated by MongoGen. DO NOT EDIT.")
	g.p("// Collection: %s", collection)
	g.p("")
	g.p("package %s", pkg.name)
	g.p("")
	g.p("import (")
	g.in()
	for path := range pkg.imports {
		g.p("%q", path)
	}
	g.out()
	g.p(")")
	g.p("")

	g.p("type %sFilter struct {", singular)
	g.in()
	g.p("Filter bson.D")
	g.out()
	g.p("}")
	g.p("")

	g.p("type %s struct {", plural)
	g.in()
	g.p("db *mongo.Database")
	g.out()
	g.p("}")
	g.p("")

	g.p("func New%s(db *mongo.Database) *%s {", plural, plural)
	g.in()
	g.p("return &%s{db}", plural)
	g.out()
	g.p("}")
	g.p("")

	g.p("func %sWithID(id primitive.ObjectID) %sFilter {", singular, singular)
	g.in()
	g.p("return %sFilter{bson.D{{Key: %q, Value: %s}}}", singular, "_id", "id")
	g.out()
	g.p("}")
	g.p("")

	g.p("func %sWithIDs(ids []primitive.ObjectID) %sFilter {", singular, singular)
	g.in()
	g.p("return %sFilter{bson.D{{Key: %q, Value: bson.M{%q: ids}}}}", singular, "_id", "$in")
	g.out()
	g.p("}")
	g.p("")

	for _, m := range pkg.methods {
		g.p("func %s(%s) %sFilter {", m.name, printMethodArgs(m.args), singular)
		g.in()
		g.p(printMethodReturn(singular, m.args))
		g.out()
		g.p("}")
		g.p("")
	}

	g.p("func (s *%s) Find(ctx context.Context, filter %sFilter, opts ...*options.FindOptions) (*mongo.Cursor, error) {", plural, singular)
	g.in()
	g.p("return s.db.Collection(%q).Find(ctx, filter.Filter, opts...)", collection)
	g.out()
	g.p("}")
	g.p("")

	g.p("func (s *%s) FindWithIDs(ctx context.Context, ids []primitive.ObjectID, opts ...*options.FindOptions) (*mongo.Cursor, error) {", plural)
	g.in()
	g.p("return s.db.Collection(%q).Find(ctx, bson.M{%q: bson.M{%q: ids}}, opts...)", collection, "_id", "$in")
	g.out()
	g.p("}")
	g.p("")

	g.p("func (s *%s) FindOne(ctx context.Context, filter %sFilter, opts ...*options.FindOneOptions) *mongo.SingleResult {", plural, singular)
	g.in()
	g.p("return s.db.Collection(%q).FindOne(ctx, filter.Filter, opts...)", collection)
	g.out()
	g.p("}")
	g.p("")

	g.p("func (s *%s) FindOneWithID(ctx context.Context, id primitive.ObjectID, opts ...*options.FindOneOptions) *mongo.SingleResult {", plural)
	g.in()
	g.p("return s.db.Collection(%q).FindOne(ctx, bson.M{%q: id}, opts...)", collection, "_id")
	g.out()
	g.p("}")
	g.p("")

	g.p("func (s *%s) Count(ctx context.Context, filter %sFilter, opts ...*options.CountOptions) (int64, error) {", plural, singular)
	g.in()
	g.p("return s.db.Collection(%q).CountDocuments(ctx, filter.Filter, opts...)", collection)
	g.out()
	g.p("}")
	g.p("")
}

func printMethodArgs(args []methodArg) string {
	out := ""
	for _, arg := range args {
		out += fmt.Sprintf("%s %s, ", arg.name, arg.typ)
	}
	return out[:len(out)-2]
}

func printMethodReturn(returnTyp string, args []methodArg) string {
	out := ""
	for _, arg := range args {
		if arg.multiple {
			out += fmt.Sprintf("{Key: %q, Value: bson.M{%q: %s}}, ", arg.query, "$in", arg.name)
		} else {
			out += fmt.Sprintf("{Key: %q, Value: %s}, ", arg.query, arg.name)
		}
	}
	return fmt.Sprintf("return %sFilter{bson.D{%s}}", returnTyp, out[:len(out)-2])
}
