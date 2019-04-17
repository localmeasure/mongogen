package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/globalsign/mgo"
	"github.com/roamz/mongogen"
)

var (
	mongoURI    = flag.String("server", "mongodb://localhost:27017", "mongo server, example: mongodb://localhost:27017")
	mongoDB     = flag.String("db", "test", "database name, example: test")
	destination = flag.String("dst", "", "output file; defaults to stdout.")
)

func main() {
	flag.Parse()
	session, err := mgo.Dial(*mongoURI)
	if err != nil {
		log.Println(err)
		return
	}

	analyzer := mongogen.NewAnalyzer(session, *mongoDB)

	pkg, err := analyzer.Analyze()
	if err != nil {
		log.Println(err)
		return
	}

	dst := os.Stdout
	if len(*destination) > 0 {
		if err := os.MkdirAll(filepath.Dir(*destination), os.ModePerm); err != nil {
			log.Fatalf("Unable to create directory: %v", err)
		}
		f, err := os.Create(*destination)
		if err != nil {
			log.Fatalf("Failed opening destination file: %v", err)
		}
		defer f.Close()
		dst = f
	}

	gen := mongogen.NewGenerator("services")
	gen.Generate(&pkg)

	if _, err := dst.Write(gen.Output()); err != nil {
		log.Fatalf("Failed writing to destination: %v", err)
	}
}