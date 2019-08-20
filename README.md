# mongogen
mongodb index-based query builder (go:generate + [mongodb-go-driver](https://github.com/mongodb/mongo-go-driver)).

### install
```
go get github.com/localmeasure/mongogen
go install github.com/localmeasure/mongogen/mongogen
```

### setup

Declare index specs:

* Singlefield index: `field1:go_type`. Example: `group_id:int`
* Compound index: `field2:go_type+field3:go_type`. Example: `location_id:id+last_seen:time`
* Multikey index: `field4:[]go_type`. Example: `team_id:[]int`

Some common types are supported:

* id: `primitive.ObjectID`
* time: `time.Time`

### examples

Create `example.go`:
```
//go:generate mongogen -p users -c users -i group_id:id+name:string -i team_id:id+last_seen:time
```

Then run (see example output [here](https://github.com/localmeasure/mongogen/blob/master/_example/mongo.go)):
```
go generate
```
