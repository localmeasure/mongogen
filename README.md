# mongogen

A codegen, generate mongo query builder in golang base on indices, group in services

## install
```
go get github.com/localmeasure/mongogen
go install github.com/localmeasure/mongogen/mongogen
```

## generate
```
mongogen -server mongodb://localhost:27017 -db localmeasure
```

## todo
* Remove globalsign/mgo (using this due to some current hard code in [mongo-go-driver](https://github.com/mongodb/mongo-go-driver))
* Slowlog unoptimized indices (e.g: objectid-date-objectid)
