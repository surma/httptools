Package httptools is a collection of simple helper types for Go’s net/http.

For details and examples, please see the [documentation](http://godoc.org/github.com/surma/httptools).

[![Build Status](https://drone.io/github.com/surma/httptools/status.png)](https://drone.io/github.com/surma/httptools/latest)

## Contrived example

```Go
r := httptools.NewRegexpSwitch(map[string]http.Handler{
	"/people/(.+)": httptools.L{
		httptools.SilentHandler(AuthenticationHandler),
		httptools.MethodSwitch{
			"GET": ListPeople,
			"PUT": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				name := strings.StripPrefix(r.URL.Path, "/people/")
				AddNewPerson(name)
			})
		},
		SaveSessionHandler,
	},
	"/.+": http.FileServer(http.Dir("./static")),
})
http.ListenAndServe("localhost:8080", r)
```

## Tools
httptools provides the following tools:
### Handler list
Define a sequence of `http.Handler`. One will be executed after another.
### Silent handler
If a silent handler produces output, it is assumed to be an error. If the
silent handler is in a handler list, the execution of that list will be aborted.
### Switches
#### Method switch
Dispatch requests to different handlers according the the HTTP verb used
in the request.
#### RegexpSwitch
Dispatch requests to different handlers according to regexps being matched
against the request path.
#### HostnameSwitch
Dispatch requests to different handlers according to the hostname used
in the request.
### Mounts
Dispatch requests to different handlers according to path prefixes. The
path prefix will be stripped from the request before being passed to the
handler.

---
Version 2.1.0
