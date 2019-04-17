# mongogen

A codegen, generate any possible IXSCAN queries in golang, group to services

first draft

* get every indexes in database
* fetch record with index key(s) to find its data type (mapping to golang)
* **generate golang code** to do safe querying
* compound index will have n (keys) of methods (in keys order)

benefits

* no need to do profiling
* no more query.explain() (generated queries will always be IXSCAN)
* cover all services, dont have to write custom service again, e.g: https://gitlab.com/localmeasure/std/tree/master/services