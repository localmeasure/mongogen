package main

import (
	"log"

	"github.com/globalsign/mgo"
	"github.com/localmeasure/mongogen"
)

func main() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	analyzer := mongogen.NewAnalyzer(session)

	pkg, err := analyzer.Analyze()
	if err != nil {
		log.Println(err)
		return
	}

	gen := mongogen.NewGenerator("services")
	gen.Generate(&pkg)
	log.Println(string(gen.Output()))
}
