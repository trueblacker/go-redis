# How to implement a Redis-compatible server?

This section describes how to implement a Redis-compatible server using the go-redis. Please see to the [examples](../examples) directory to know how to implement the compatible server in more detail.

## Creating your Redis-compatible server

### STEP1: Inheritancing the base server

At first, inherit the base server of the go-redis to implement your Redis-compatible server as the following:

```
import (
	"github.com/trueblacker/go-redis/redis"
)

type Server struct {
	*redis.Server
}
```

### STEP2: Implementing your user command handler

Next, implement your user command handler such as SET and GET according to the [UserCommandHandler](../redis/handler.go) interface of the go-redis as the following:

```
func (server *Server) Set(conn *redis.Conn, key string, val string, opt redis.SetOption) (*redis.Message, error) {
	dbID := conn.Database()

    ....

	return redis.NewOKMessage(), nil
}

func (server *Server) Get(conn *redis.Conn, key string) (*redis.Message, error) {
	dbID := conn.Database()

    ....

	return redis.NewStringMessage(string(record.Data)), nil
}
```

The Conn has the connection information such as the selected database identifier, and the all handler methods should return the appropriate RESP message response.

### STEP4: Setting your user command handler

Next, set your user command handler to your server using `Server::SetCommandHandler()` as the following:

```
func NewServer() *Server {
	server: = &Server{
		Server:    redis.NewServer(),
	}
    server.SetCommandHandler(server)
    return server
}
```

### STEP5: Starting server

Finally, start your compatible server using `Server::Start()` as the following:

```
server := NewServer()
....
server::Start()
```

## Next Steps

To know the server implementation using go-redis in more detail, let's check the following documentation 

- [Inside of go-redis](server_inside.md)
- Examples
    - [go-redisd](../examples/go-redisd)
