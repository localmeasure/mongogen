# mongogen

A codegen, generate any possible IXSCAN queries in golang, group to services

first draft

* get every collection names
* get every indexes
* fetch record with index key(s) exist to find its data type to map to golang
* **generate golang code** to do safe querying
* compound index will have n (keys) of methods (in keys order)

benefits

* no need to do profiling
* no more query.explain() (generated queries will always be IXSCAN)
* cover all services, dont have to write custom service again, e.g: https://gitlab.com/localmeasure/std/tree/master/services