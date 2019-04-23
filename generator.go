package mongogen

import (
	"bytes"
	"fmt"
	"strings"
)

type Generator struct {
	PackageName string
	buf         bytes.Buffer
	indent      string
}

func NewGenerator(packageName string) *Generator {
	return &Generator{PackageName: packageName}
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

func (g *Generator) Generate(pkg *Pkg) {
	g.p("// Code generated by MongoGen. DO NOT EDIT.")
	g.p("// Database: %v", pkg.Database)
	g.p("// \t+---------------------------+---------------------------+")
	g.p("// \t| Collection                | Type                      |")
	g.p("// \t+---------------------------+---------------------------+")
	for _, t := range pkg.Types {
		g.p("// \t" + t.Collection + strings.Repeat(" ", 30-len(t.Collection)) + t.Plural)
	}
	g.p("")
	g.p("package %v", g.PackageName)
	g.p("")
	g.p("import (")
	g.in()
	for _, path := range pkg.Imports {
		g.p("%q", path)
	}
	g.out()
	g.p(")")
	g.p("")

	g.p("const (")
	g.in()
	for _, constant := range pkg.Consts {
		g.p("%v "+strings.Repeat(" ", 30-len(constant[0]))+"= %q", constant[0], constant[1])
	}
	g.out()
	g.p(")")
	g.p("")

	g.p("type Filter struct {")
	g.in()
	g.p("Filter bson.D")
	g.p("Hint   string")
	g.out()
	g.p("}")
	g.p("")

	for _, t := range pkg.Types {
		g.p("type %v struct {", t.Plural)
		g.in()
		g.p("db *mongo.Database")
		g.out()
		g.p("}")
		g.p("")

		g.p("func New" + t.Plural + "(db *mongo.Database) *" + t.Plural + " {")
		g.in()
		g.p("return &" + t.Plural + "{db}")
		g.out()
		g.p("}")
		g.p("")

		for _, m := range t.Methods {
			g.p("func " + m.Name + "(" + printMethodArgs(m.Args) + ") Filter {")
			g.in()
			g.p(printMethodReturn(m.Args, m.Hint))
			g.out()
			g.p("}")
			g.p("")
		}

		g.p("func (s *" + t.Plural + ") Find(ctx context.Context, filter Filter, opts ...*options.FindOptions) (*mongo.Cursor, error) {")
		g.in()
		g.p("opts = append(opts, options.Find().SetHint(filter.Hint))")
		g.p("return s.db.Collection(%q).Find(ctx, filter.Filter, opts...)", t.Collection)
		g.out()
		g.p("}")
		g.p("")

		g.p("func (s *" + t.Plural + ") FindWithIDs(ctx context.Context, ids []primitive.ObjectID, opts ...*options.FindOptions) (*mongo.Cursor, error) {")
		g.in()
		g.p("return s.db.Collection(%q).Find(ctx, bson.M{%q: bson.M{%q: ids}}, opts...)", t.Collection, "_id", "$in")
		g.out()
		g.p("}")
		g.p("")

		g.p("func (s *" + t.Plural + ") FindOne(ctx context.Context, filter Filter, opts ...*options.FindOneOptions) *mongo.SingleResult {")
		g.in()
		g.p("opts = append(opts, options.FindOne().SetHint(filter.Hint))")
		g.p("return s.db.Collection(%q).FindOne(ctx, filter.Filter, opts...)", t.Collection)
		g.out()
		g.p("}")
		g.p("")

		g.p("func (s *" + t.Plural + ") FindOneWithID(ctx context.Context, id primitive.ObjectID, opts ...*options.FindOneOptions) *mongo.SingleResult {")
		g.in()
		g.p("return s.db.Collection(%q).FindOne(ctx, bson.M{%q: id}, opts...)", t.Collection, "_id")
		g.out()
		g.p("}")
		g.p("")
	}
}

func printMethodArgs(args []Argument) string {
	out := ""
	for _, arg := range args {
		out += arg.ArgName + " " + arg.ArgType + ", "
	}
	return out[:len(out)-2]
}

func printMethodReturn(args []Argument, hint string) string {
	out := ""
	for _, arg := range args {
		out += "{Key: \"" + arg.QueryName + "\", Value: " + arg.ArgName + "}, "
	}
	return "return Filter{bson.D{" + out[:len(out)-2] + "}, " + hint + "}"
}

func (g *Generator) Output() []byte {
	return g.buf.Bytes()
}
