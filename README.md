# mongogen

A codegen, generate mongo query builder in golang base on indices, group in services

## install
```
go get github.com/roamz/mongogen
go install github.com/roamz/mongogen/mongogen
```

## generate
```
mongogen -server mongodb://localhost:27017 -db localmeasure
```

Output example:
https://gitlab.com/localmeasure/std/blob/master/services/mongo.go

TODO:
* Remove globalsign/mgo (using this due to some current hard code in [mongo-go-driver](https://github.com/mongodb/mongo-go-driver))
* Slowlog unoptimized indices (e.g: objectid-date-objectid)
