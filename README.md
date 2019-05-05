# mongogen
static codegen for mongo queries (go:generate)

### install
```
go get github.com/localmeasure/mongogen
go install github.com/localmeasure/mongogen/mongogen
```

### example

Declare mongo index specs in .go file
```
//go:generate mongogen -p users -c users -i group_id:id+name:string -i team_id:id+last_seen:time
```

Then run
```
go generate
```
See sample output [here](https://github.com/localmeasure/mongogen/blob/master/_example/mongo.go).
