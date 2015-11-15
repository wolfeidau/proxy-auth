# proxy-auth

This is a small module which provides a small authentication interface you can embed into proxy servers. It currently supports authentication to GitHub using Oauth2.

# usage

Download this package.

```
go get github.com/wolfeidau/proxy-auth
```

An example server is shown below.

```go
package main

import (
    "net/http"

    "github.com/gorilla/mux"
    "github.com/gorilla/sessions"
    "github.com/wolfeidau/proxy-auth"
)

func main() {
    // setup a store, in our case one using secure cookies
    store := sessions.NewCookieStore([]byte("something-very-secret"))
    s := auth.NewServer(store)

    // configure a mux
    r := mux.NewRouter()
    r.PathPrefix(auth.PathPrefix).Handler(s.GetMux())

    // add a wrapper to check the session for each request
    o := auth.CheckSession(r, store)

    // listen to the network
    http.ListenAndServe(":5000", o)
}
```


# License

This code is Copyright (c) 2014 Mark Wolfe and licensed under the MIT license. All rights not explicitly granted in the MIT license are reserved. See the included `LICENSE.txt` file for more details.
