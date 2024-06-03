MODULE = github.com/xiaoyi-byte/student-demo

# generate client code by IDL.
.PHONY: gen
gen:
	kitex -module $(MODULE) -service student idl/student.thrift

# run the gateway
.PHONY: gateway
gateway:
	go run .
	go run hertz-gateway/main.go

