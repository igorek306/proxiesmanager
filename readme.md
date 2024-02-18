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

#### Initialization
```go
m := proxiesmanager.NewProxiesManager()
```
#### Configuration
```go
/*
when the manager didn't find any "://" in the line when loading the proxy, it assumes there is no schema and adds a default scheme
example:
127.0.0.1:5555 -> http://127.0.0.1:5555
*/
m.DefaultScheme("http") 
/*
when using pm.CallRequest or pm.LoadFromWeb with proxies,manager tries call request using each loaded proxy and the tries value determines how many times it will test the same proxy before returning an error
*/
m.DefaultTryAttempts(2)
/*
these two functions return the current value, so if you want to check the default values without changing them, try:
*/
fmt.Printf("Default scheme is: %s, Default attempts count: %d", m.DefaultScheme(""),m.DefaultTryAttempts(0))
```