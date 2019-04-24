package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/globalsign/mgo"
	"github.com/localmeasure/mongogen"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoURI    = flag.String("server", "mongodb://localhost:27017", "mongo server, example: mongodb://localhost:27017")
	mongoDB     = flag.String("db", "test", "database name, example: test")
	destination = flag.String("dst", "", "output file; defaults to stdout.")
	packageName = flag.String("pkg", "services", "package name, example: services")
)

func main() {
	flag.Parse()
	session, err := mgo.Dial(*mongoURI)
	if err != nil {
		log.Println(err)
		return
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(*mongoURI))
	if err != nil {
		log.Println(err)
		return
	}

	analyzer := mongogen.NewAnalyzer(session, client.Database(*mongoDB), *mongoDB)
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

	gen := mongogen.NewGenerator(*packageName)
	gen.Generate(&pkg)
	if _, err := dst.Write(gen.Output()); err != nil {
		log.Fatalf("Failed writing to destination: %v", err)
	}
}
