module github.com/pyihe/go-example/grpc

go 1.18

require (
	github.com/pyihe/go-pkg v0.0.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.1
)

require (
	cloud.google.com/go v0.34.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20220812174116-3211cb980234 // indirect
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.4.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)

replace github.com/pyihe/go-pkg v0.0.0 => ../../go-pkg
