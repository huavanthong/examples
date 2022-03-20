# Introduction
This sample will help you know how do seperate your service to each layer, it help you easy to manage it.
More explain: [here](https://gokit.io/examples/stringsvc.html#middlewares)

# Question
## About Separation of concerns knowledge
* [If you know about example: stringsvc1, do you think that it's hard complexity for maintenance? ](https://github.com/huavanthong/examples/tree/master/stringsvc1)
* [Why we need to seperate packge into smaller part?](#why-we-need-to-seperate-package)

## About Transport logging and Application logging
* [Why we need to use log in our service? Where could we put log in scope in our program?](#why-we-need-logging)
* [How do you implement transport logging](#transport-logging)
* [How do you implment application logging? What the difference between them?](#application-logging)
* [Could you tell me how many parameters we need to pass to logger?](#parameters-for-logger)
## About Application instrumentation
* [What is instrumentation?](#https://www.meisternote.com/app/note/rvWxC3RKXtQm/application-instrumentation)
* [How do you metric in this example?](#to-get-metris)


## Answer question
### Why we need to seperate package?
Separating each layer of the call graph into individual files makes a go-kit project easier to read as you increase the number of service endpoints.
- The first reason, it is that we can maintenance it easyly.
- The second reason, we can add more services to our project.
- The third reason, we can manage endponts, and transaction in our server.


### Why we need logging
Any component that needs to log should treat the logger like a dependency, same as a database connection. So, we construct our logger in our func main, and pass it to components that need it.
- We never use a globally-scoped logger.
===> Why we shouldn't put log at the globally-scoped logger

### Transport logging
You can follow steps to implement transport logging
1. TransportLoggingS1: Create format logger
2. TransportLoggingS2: create a object contain our service
3. TransportLoggingS3: Use chain handler to make logger


### Parameters for logger
One of the keyword, you need to remember that are what parameters you need to pass to logger. It includes:
- What method you are using?
- What input for this method?
- What output for this method? 
- What error on this method if it occurs?
- This method take how much time for running?
Through that, we summary as below:
```
    mw.logger.Log(
                "method", "uppercase",
                "input", s,
                "output", output,
                "err", err,
                "took", time.Since(begin),
            )
```
### Steps to implement Application Logger
After you're passing logger through transport, you continue to pass logger to your application.
4. TransportLoggingS4: Create handler for Uppercase
5. TransportLoggingS5: Inside that handler, we get handler out and use it.

```
type loggingMiddleware struct {
	logger log.Logger
	next   StringService
}

func (mw loggingMiddleware) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.Uppercase(s)
	return
}
```
More details: [here](https://github.com/huavanthong/examples/blob/master/stringsvc2/logging.go)
### Application instrumentation


# Getting Started
To this project
```
> go build
```

To start this server
```
>.\stringsvc2.exe
```

To run uppercase API
```
> curl -X POST -d "{\"s\":\"hello, world\"}" localhost:8080/uppercase

method=uppercase input="hello, world" output="HELLO, WORLD" err=null took=0s   <===== logging
{"v":"HELLO, WORLD"}
```

To run count API
```
> curl -X POST -d "{\"s\":\"hello, world\"}" localhost:8080/count 
method=count input="hello, world" n=12 took=0s                                  <===== logging
{"v":12}
```
### To get metris
```
    http://localhost:8080/metrics
```
Output
```
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0
go_gc_duration_seconds_sum 0
go_gc_duration_seconds_count 0
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 7
```