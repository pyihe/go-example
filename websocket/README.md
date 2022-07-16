### Websocket 客户端-服务器

基于[github.com/gorilla/websocket](https://github.com/gorilla/websocket)实现的websocket客户端-服务器, 为每个Conn开启两个Goroutine(读写协程), 实例代码参考[echo](https://github.com/pyihe/go-example/tree/master/websocket/echo)