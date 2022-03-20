# Introduction
This sample will help you how do begin to start implement a service with go-kite framework.
More explain: [here](https://gokit.io/examples/stringsvc.html#first-principles)

# Question
## About inteface knowledge
* [Do you understand about Interface in Golang?](https://github.com/huavanthong/MasterGolang/tree/main/01_GettingStarted)
* [Could see the difference between "type StringService interface" and "type stringService struct"? Why we need stringServer?]()

## About structure data 
* [Do you know how to use json marshalling for request/response?]()
* [Carefully, we need to think about data structure in every request and response]()
## About Microservice knowledge
* [What is Endpoints in microservice?]()
* [How do you implement a Endpoint using go-kite?]()
* [What is httptransport in go-kite? How to use it?]()
* [Why we need to use context in every decode/encode a request?]()


# Getting Started
To start this server
```
    go run .\main.go
```

To run uppercase API
```
> curl -XPOST -d'{"s":"hello, world"}' localhost:8080/uppercase

{"v":"HELLO, WORLD"}
```

To run count API
```
> curl -XPOST -d'{"s":"hello, world"}' localhost:8080/count
{"v":12}
```

