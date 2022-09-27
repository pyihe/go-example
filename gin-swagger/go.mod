module github.com/pyihe/go-example/gin-swagger

go 1.18

require (
	github.com/pyihe/go-example v0.0.0-20220926144557-6ae9f6da4d95 // indirect
	github.com/pyihe/go-pkg v0.0.0-20220816061532-b61575b24296 // indirect
	github.com/pyihe/plogs v1.0.1 // indirect
	pkg v0.0.0
)

replace (
	pkg v0.0.0 => ../pkg
	github.com/pyihe/go-pkg v0.0.0-20220816061532-b61575b24296 => ../../go-pkg
)
