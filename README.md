# Hprose for nats

## support
+ timeout
+ failround

## client
```golang
import openRpc from "hprose-react";
client := rpc.NewClient("nats://localhost:4222?topic=cnmb")
	method := &struct{ Hello func(string) (string, error) }{}
	client.UseService(method)
	for i := 0; i < 3000; i++ {
		log.Print(method.Hello("baby"))
	}
});
```

## server
```golang
server := rpc.NewNatsServer(rpc.Option(rpc.Uri("nats://localhost:4222?topic=cnmb")))
	server.AddFunction("hello", func(msg string) string {
		log.Print("hello: ", msg)
		return "hi bitch!"
	})
	server.Start()
});
```