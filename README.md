# mongogen
static codegen for mongo queries (go:generate)

### install
```
go install github.com/localmeasure/mongogen/mongogen
```

### example

Declare mongo index specs in .go file
```
//go:generate mongogen -c users -i group_id:id+name:string -i team_id:id+last_seen:time
```

Then run
```
go generate
```
See sample output [here](https://github.com/localmeasure/mongogen/blob/master/_example/mongo_users.go).
