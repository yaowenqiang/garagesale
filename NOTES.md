> go mod init github.com/yaowenqiang/garagesale
$ hey http://localhost:8888

> https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html
> https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/

# pprof
> hey -c 10 -n 15000 http://localhost:8000/v1/products
> go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
> (pprof) top
> (pprof) top -cum
> (pprof) web

> go get github.com/divan/expvarmon
> 
## kubernetes


> brew install kind
> brew install kustomize

> docker pull kindest/node:v1.20.0
> kubectl version --client

> docker pull postgres:13-alpine
> docker pull openzipkin/zipkin:2.23
> docker pull alpine:3.12.3
