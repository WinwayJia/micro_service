all: consul-register consul-health

consul-register: register.go
	go build -o $@ $^

consul-health: health.go
	go build -o $@ $^


