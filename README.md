# mongogen
static codegen  for mongodb index queries (go:generate), code generated uses official mongodb go driver (https://github.com/mongodb/mongo-go-driver)

### install
```
go get github.com/localmeasure/mongogen
go install github.com/localmeasure/mongogen/mongogen
```

### example

Declare mongodb index specs in a `.go` file with `-i`, `-p` for destination go package, `-c` for mongodb collection
```
//go:generate mongogen -p users -c users -i group_id:id+name:string -i team_id:id+last_seen:time
```

Then run
```
go generate
```
See sample output [here](https://github.com/localmeasure/mongogen/blob/master/_example/mongo.go).
