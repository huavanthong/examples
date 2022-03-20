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
* [How do you convert any request from client into our structure data?](#convert-request-to-our-struct-data)
## About Microservice knowledge
* [What is Endpoints in microservice?]()
* [How do you implement a Endpoint using go-kite?]()
* [What process inside a Endpoint?](#convert-request-to-our-struct-data)
* [What is httptransport in go-kite? How to use it?]()
* [Why we need to use context in every decode/encode a request?]()

## Questions
### Convert request to our struct data
```
func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		// Step 5.1: Convert any request interface to our structure.
		req := request.(uppercaseRequest)

		// Step 5.2: Use our service to upper case data from request.
		v, err := svc.Uppercase(req.S)

		// Step 5.3: Response output
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}
```

### Implement Endpoint with go-kite
Step 1: Implement structure data
```
type StringService interface {
	Uppercase(string) (string, error)
    .....
}
```
Step 2: Create a Endpoint
```
func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint  {
    ........
}
```
# Getting Started
To start this server
```
    go run .\main.go
```

To run uppercase API
```
> curl -X POST -d "{\"s\":\"hello, world\"}" localhost:8080/uppercase

{"v":"HELLO, WORLD"}
```

To run count API
```
> curl -X POST -d "{\"s\":\"hello, world\"}" localhost:8080/count
{"v":12}
```
