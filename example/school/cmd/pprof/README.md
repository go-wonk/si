# Pprof

## Profile
### CPU
```bash
go tool pprof http://127.0.0.1:6060/debug/pprof/profile
```

### Heap
```bash
go tool pprof http://127.0.0.1:6060/debug/pprof/heap
```

## View on web
```
go tool pprof -http=:9999 /path/to/pprof/pprof.samples.cpu.001.pb.gz
```

## Trace
wget http://127.0.0.1:6060/debug/pprof/trace?seconds=10 -O /path/to/trace
go tool trace /path/to/trace


## Bench test
### bombardier
bombardier -c 20 -n 1000 \
-H "Content-Type: application/json; charset=utf-8" \
-m POST -k \
-b '{"name":"wonk","email_address":"wonk@wonk.org"}' \
http://172.16.200.82:8080/test/pprof

bombardier -c 20 -n 1000 \
-H "Content-Type: application/json; charset=utf-8" \
-m POST -k \
-b '{"name":"wonk","email_address":"wonk@wonk.org"}' \
http://172.16.200.82:8080/test/findall

### gc
curl --insecure -H "Content-Type: application/json; charset=utf-8" \
-X GET \
http://127.0.0.1:8080/test/gc

## Curl
curl --insecure -H "Content-Type: application/json; charset=utf-8" \
-X POST \
-d '{"name":"wonk","email_address":"wonk@wonk.org"}' \
http://127.0.0.1:8080/test/findall

curl --insecure -H "Content-Type: application/json; charset=utf-8" \
-X POST \
-d '{"name":"wonk","email_address":"wonk@wonk.org"}' \
http://127.0.0.1:8080/test/repeat/findall

