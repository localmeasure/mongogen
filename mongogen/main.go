// Copyright 2019 Local Measure. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/localmeasure/mongogen"
)

var collection = flag.String("c", "", "-c users")

type indexes []string

func (i *indexes) String() string {
	return fmt.Sprint(*i)
}

func (i *indexes) Set(value string) error {
	for _, idx := range strings.Split(value, ",") {
		*i = append(*i, idx)
	}
	return nil
}

var indexFlags indexes

func init() {
	flag.Var(&indexFlags, "i", "-i group_id:id+name:string -i team_id:id+last_seen:time")
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	flag.Parse()
	g := mongogen.NewGenerator()
	g.Gen(*collection, indexFlags)
	f, err := os.Create("mongo_" + strings.ToLower(*collection) + ".go")
	if err != nil {
		log.Fatalf("Failed creating file: %v", err)
	}
	defer f.Close()
	if _, err := f.Write(g.Output()); err != nil {
		log.Fatalf("Failed writing to destination: %v", err)
	}
}
