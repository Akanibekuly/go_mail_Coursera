curl http://127.0.0.1:8080/debug/pprof/heap -o mem_out.txt
curl http://127.0.0.1:8080/debug/pprof/profile?seconds=5 -o cpu_out.txt

go tool pprof -svg -alloc_objects mem_out.txt > mem_ao.svg
go tool pprof -svg cpu_out.txt > cpu.svg