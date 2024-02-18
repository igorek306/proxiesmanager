# Proxies Manager for Go
 this is small set of functions that can simplify your work with proxies  
 main target is to use that with http requests, take a look at examples

## Installation 
```bash
go get github.com/igorek306/proxiesmanager
```
then in your Go file
```go
import "github.com/igorek306/proxiesmanager"
```
## Usage

### Initialization
```go
m := proxiesmanager.NewProxiesManager()
```
### Configuration
```go
/*
when the manager didn't find any "://" in the line when loading the proxy,
it assumes there is no schema and adds a default scheme, example:
127.0.0.1:5555 -> http://127.0.0.1:5555
*/
m.DefaultScheme("http") 
/*
when using pm.CallRequest or pm.LoadFromWeb with proxies,
manager tries call request using each loaded proxy and the tries value
determines how many times it will test the same proxy before returning an error
*/
m.DefaultTryAttempts(2)
/*
these two functions return the current value,
so if you want to check the default values without changing them, try:
*/
fmt.Printf("Default scheme is: %s, Default attempts count: %d", m.DefaultScheme(""),m.DefaultTryAttempts(0))
```
### Loading proxies from file
```go
n, err := m.LoadFromFile("proxies.txt")
if err != nil {
    fmt.Printf("error loading proxies, %s\n", err.Error())
    return
}
fmt.Printf("loaded %d proxies from file\n", n)
```
### Loading proxies from web
```go
/*
this example will do a GET request for the given URL
and load proxies like from a file, but from the response body
*/
n, err := m.LoadFromWeb(proxiesmanager.TargetSite{
    Url: "https://leaked.wiki/r/JFNAXhn1oS",
}, nil, false)
if err != nil {
    fmt.Printf("error loading proxies, %s\n", err.Error())
    return
}
fmt.Printf("loaded %d proxies from web\n", n)
```
```go
/*
this example does the same as above, but does the request
using already loaded proxies (if none loaded returns error)
*/
n, err := m.LoadFromWeb(proxiesmanager.TargetSite{
    Url: "https://leaked.wiki/r/JFNAXhn1oS",
}, nil, true)
```
```go
/* 
this example loads proxies using given request
useful when, for example, website requires certain headers 
*/
req, _ := http.NewRequest("GET", "https://leaked.wiki/r/JFNAXhn1oS", nil)
req.Header.Add("User-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")
n, err := m.LoadFromWeb(proxiesmanager.TargetSite{
    Request: req,
}, nil, false)
```
```go
/*
you can pass prepared http.Client too, like:
*/
myHttpClient := &http.Client{}
// configure client as needed
n, err := m.LoadFromWeb(proxiesmanager.TargetSite{
    Url: "https://leaked.wiki/r/JFNAXhn1oS",
}, myHttpClient, false)
```
### Add proxies manually
```go
proxy1, _ := url.Parse("http://example.com:500")
proxy2, _ := url.Parse("http://example.com:1000")
m.InsertProxies(&[]*url.URL{
    proxy1, proxy2,
})
```
### Calling http requests using proxies
```go
// preparing request
req, _ := http.NewRequest(
    "GET",
    "https://httpbin.org/ip",
    nil,
)
```
```go
// calling with default client
res, err := m.CallRequest(
    nil, req,
)
```
or
```go
// calling with custom client
myHttpClient := &http.Client{}
res, err := m.CallRequest(
    myHttpClient, req,
)
```
```go
// handling
if err != nil {
    fmt.Printf("error doing request, %s\n", err)
    return
}
defer res.Body.Close()
fmt.Printf("Got response using %s proxy,  body:\n\n", m.Proxy().String())
io.Copy(os.Stdout, res.Body)
```
### Other functions
#### Some of them have already been used in examples. 
#### These are useful if you want to implement proxy servers yourself.
```go
// checks the number of proxies already loaded
fmt.Printf("Already loaded %d proxies", m.Count())
```
```go
// returns proxy to use
fmt.Printf("Successfully called request using %s proxy", m.Proxy())
// default current proxy is index = 0
// so for example, first line of loaded file
```
```go
// to make manager selects next proxy as the current one, use
m.Next()
```
```go
// if you want to call m.next and then m.proxy,
// use this shorthand:
fmt.Printf("Now requesting with %s proxy", m.NextProxy())
```
```go
// if you rerun your app quite frequently,
// it's recommended to set index to random
// to prevent using the same proxy, to the same action
m.SelectRandom()
```
```go
m.Proxies() // returns pointer to the manager's proxies array
```
```go
// don't use above to for example, write all loaded proxies to file
// use this func instead:
m.PrintAll(file) 
// or if you want to print all proxies:
m.PrintAll(os.Stdout)
```