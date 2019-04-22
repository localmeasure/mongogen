package mongogen

import (
	"context"
	"log"
	"runtime"
	"sync"

	"github.com/globalsign/mgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	bsonTyps = [...]int{7, 2, 16, 18, 8, 1, 9, 3, 4, 17, 5}
	typMaps  = map[int]string{
		7:  "primitive.ObjectID",
		2:  "string",
		16: "int",
		18: "int", // should be int64 on 32-bit arch
		8:  "bool",
		1:  "float64",
		9:  "time.Time",
		3:  "primitive.M",
		4:  "primitive.A",
		17: "primitive.Timestamp",
		5:  "primitive.Binary",
		// bson types below will be mapped to interface{}
		// 6:  "primitive.Undefined",
		// 10: "primitive.Null",
		// 11: "primitive.Regex",
		// 12: "primitive.DBPointer",
		// 13: "primitive.JavaScript",
		// 14: "primitive.Symbol",
		// 15: "primitive.CodeWithScope",
		// 19: "primitive.Decimal128",
	}
	setupWorkers sync.Once
	typReg       struct {
		sync.RWMutex
		ref map[string]string
	}
)

type typCheck struct {
	Key        string
	Query      bson.D
	Hint       string
	Typ        string
	BsonTyp    int
	Collection string
}

func (a *Analyzer) typesOf(colName string, idx mgo.Index) {
	a.startWorkers()
	var query bson.D
	for i := 0; i < len(idx.Key); i++ {
		// query is an ordered bson doc
		query = append(query, primitive.E{Key: idx.Key[i]})
		go func(query bson.D, i int) {
			for _, n := range bsonTyps {
				cp := make(bson.D, len(query))
				copy(cp, query)
				cp[i].Value = bson.M{"$type": n}
				a.checkc <- typCheck{
					Key:        idx.Key[i],
					Query:      cp,
					Hint:       idx.Name,
					Collection: colName,
					BsonTyp:    n,
				}
			}
		}(query, i)
		for range bsonTyps {
			tc := <-a.donec
			if tc.Typ != "" {
				typReg.Lock()
				typReg.ref[tc.Collection+tc.Key] = typMaps[tc.BsonTyp]
				typReg.Unlock()
				// set $type for index prefix, next prefix will be filter based on this
				query[i].Value = bson.M{"$type": tc.BsonTyp}
				break
			}
		}
		// unable to find supported $type for index prefix, abort
		if query[i].Value == nil {
			break
		}
	}
}

func (a *Analyzer) startWorkers() {
	setupWorkers.Do(func() {
		n := runtime.NumCPU()
		if n < 4 {
			n = 4
		}
		a.checkc = make(chan typCheck, n)
		a.donec = make(chan typCheck, n)
		typReg.ref = make(map[string]string)
		for i := 0; i < n; i++ {
			go func(in <-chan typCheck, out chan<- typCheck) {
				ctx := context.Background()
				for {
					tc := <-in
					cur, err := a.db.Collection(tc.Collection).Find(ctx, tc.Query, options.Find().SetHint(tc.Hint))
					if err != nil {
						out <- tc
						log.Println(err)
						continue
					}
					if cur.Next(ctx) {
						tc.Typ = typMaps[tc.BsonTyp]
					}
					out <- tc
				}
			}(a.checkc, a.donec)
		}
	})
}
