# Introduction
This sample will help you know how do use proxy in your server to hide our service.
More explain: [here](https://gokit.io/examples/stringsvc.html#calling-other-services)

# Question
## About Separation of concerns knowledge

## About Service discovery and load balancing
* [How do you set a optional port for proxy from input user?](#optional-listion-port-and-proxy)
* [How do you implement a proxy for load balance?](#circuite-breaking---load-balance)
* [How do you server discovery your service?](#service-discovery)


## Answer question
### Optional listion port and proxy
To set listen port from input user
```
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
		proxy  = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy uppercase requests")
	)
	flag.Parse()
```

To run it
```
>.\stringsvc3.exe -listen=:8001
```
### Circuite Breaking - Load balance
**Step 1:** Input proxy from user
```
	svc = proxyingMiddleware(context.Background(), *proxy, logger)(svc)
```
**Step 2:** Check instance of proxy from user
```
func proxyingMiddleware(ctx context.Context, instances string, logger log.Logger) ServiceMiddleware {

	if instances == "" {
		logger.Log("proxy_to", "none")
		return func(next StringService) StringService { return next }
	}

    ...
```

**Step 3:** Create Backoff for each request, and setting max time each request.
```
	var (
		qps         = 100                    // beyond which we will return an error
		maxAttempts = 3                      // per request, before giving up
		maxTime     = 250 * time.Millisecond // wallclock time, before giving up
	)
```

**Step 4:** Seperate instance to get multiple proxy from user, then we implement circuit breaking for our server.
**Input:**
```
>.\stringsvc3.exe -listen=:8080 -proxy=localhost:8001,localhost:8002,localhost:8003
```

**Implement:**
```

	var (
		instanceList = split(instances)
		endpointer   sd.FixedEndpointer
	)

	logger.Log("proxy_to", fmt.Sprint(instanceList))
	for _, instance := range instanceList {
		var e endpoint.Endpoint
		e = makeUppercaseProxy(ctx, instance)
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		endpointer = append(endpointer, e)
	}


func split(s string) []string {
	a := strings.Split(s, ",")
	for i := range a {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}

```

**Output:**
```
listen=:8001 caller=logging.go:28 method=uppercase input=foo output=FOO err=null took=5.168µs
listen=:8080 caller=logging.go:28 method=uppercase input=foo output=FOO err=null took=4.39012ms
listen=:8002 caller=logging.go:28 method=uppercase input=bar output=BAR err=null took=5.445µs
listen=:8080 caller=logging.go:28 method=uppercase input=bar output=BAR err=null took=2.04831ms
listen=:8003 caller=logging.go:28 method=uppercase input=baz output=BAZ err=null took=3.285µs
listen=:8080 caller=logging.go:28 method=uppercase input=baz output=BAZ err=null took=1.388155ms
```

**Step 5:** Implement load balance for our server
```
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)

	return func(next StringService) StringService {
		return proxymw{ctx, next, retry}
	}
}
```

### Service discovery
**Step 1:** Define interface from proxy to refer next endpoint
```
// proxymw implements StringService, forwarding Uppercase requests to the
// provided endpoint, and serving all other (i.e. Count) requests via the
// next StringService.
type proxymw struct {
	ctx       context.Context
	next      StringService     // Serve most requests via this service...
	uppercase endpoint.Endpoint // ...except Uppercase, which gets served by this endpoint
}
```

**Step 2:** Define handler for proxy discovery
```
func (mw proxymw) Uppercase(s string) (string, error) {
	response, err := mw.uppercase(mw.ctx, uppercaseRequest{S: s})
	if err != nil {
		return "", err
	}

	resp := response.(uppercaseResponse)
	if resp.Err != "" {
		return resp.V, errors.New(resp.Err)
	}
	return resp.V, nil
}
```

**Step 3:** Make uppercase proxy
```
func makeUppercaseProxy(ctx context.Context, instance string) endpoint.Endpoint {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}
	if u.Path == "" {
		u.Path = "/uppercase"
	}
	return httptransport.NewClient(
		"GET",
		u,
		encodeRequest,
		decodeUppercaseResponse,
	).Endpoint()
}

```

# Getting Started
## Build
To build this server
```
go build
```

## Run server
To run this server, with default port 8080
```
>.\stringsv3.exe
listen=:8080 caller=proxying.go:27 proxy_to=none
listen=:8080 caller=main.go:68 msg=HTTP addr=:8080
```

To run this server, with another port 8001
```
>.\stringsvc3.exe -listen=:8001
isten=:8001 caller=proxying.go:27 proxy_to=none
listen=:8001 caller=main.go:68 msg=HTTP addr=:8001


>.\stringsvc3.exe -listen=:8002
listen=:8002 caller=proxying.go:27 proxy_to=none
listen=:8002 caller=main.go:68 msg=HTTP addr=:8002


>.\stringsvc3.exe -listen=:8003
listen=:8003 caller=proxying.go:27 proxy_to=none
listen=:8003 caller=main.go:68 msg=HTTP addr=:8003
```

To run this server with multiple port
```
>.\stringsvc3.exe -listen=:8080 -proxy=localhost:8001,localhost:8002,localhost:8003
listen=:8080 caller=proxying.go:46 proxy_to="[localhost:8001 localhost:8002 localhost:8003]"
listen=:8080 caller=main.go:68 msg=HTTP addr=:8080
```

## Execute
To loop 3 value
```
for s in foo bar baz ; do curl -d "{\"s\":\"$s\"}" localhost:8080/uppercase ; done
```

To loop mutliple time
```
for s in 1 2 3 4 5 6 7 8 foo bar baz ; do curl -d "{\"s\":\"$s\"}" localhost:8080/uppercase ; done
```

Infinite loop
```
for (( ; ; ))
do
    for s in foo bar baz ; do curl -d "{\"s\":\"$s\"}" localhost:8080/uppercase ; done
    #echo "infinite loops [ hit CTRL+C to stop]"
done
```