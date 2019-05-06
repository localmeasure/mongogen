# mongogen
static codegen for mongodb index queries (go:generate), generated code uses official [mongodb-go-driver](https://github.com/mongodb/mongo-go-driver).

### install
```
go get github.com/localmeasure/mongogen
go install github.com/localmeasure/mongogen/mongogen
```

### example

Mongodb index specs are denoted like below:

* Singlefield index: `field1:go_type`. Example: `group_id:int`
* Compound index: `field2:go_type+field3:go_type`. Example: `location_id:id+last_seen:time`
* Multikey index: `field4:[]go_type`. Example: `team_id:[]int`

As you notice above, other than go basic types, this codegen also supports some common types below:
* id: `primitive.ObjectID`
* time: `time.Time`

Declare mongodb index specs in a `.go` file, an index spec starts with `-i`,  `-p` for destination package name, `-c` for mongodb collection. Full example:
```
//go:generate mongogen -p users -c users -i group_id:id+name:string -i team_id:id+last_seen:time
```

Then run
```
go generate
```

See sample output [here](https://github.com/localmeasure/mongogen/blob/master/_example/mongo.go).
