# Gommon toolkit

Gommon is a set of packages for quickly prototyping Go applications.
**This stuff is very new and subject to lots of iteration. Feel free to look
and riff but don't expect any sort of stability.**


#### Auth

Auth is a simple and dumb way to store and authenticate users. It builds on the
Gorilla toolkit's sessions package (http://www.gorillatoolkit.org/pkg/sessions)
and uses SQLite to store user information. Please note, passwords are currently
stored in a very insecure fashion.


#### Render

Render is a collection of basic functions that help render template output.
If the request is an XHR request it will return a JSON response instead
of rendering it's given template. You can also append `?json` to any URL to
force a JSON response.


#### Spokes

Spokes is a basic WebSocket pub/sub implemetation using the Gorilla WebSocket
toolkit. It's based off the [chat example](https://github.com/gorilla/websocket/tree/master/examples/chat)
but adds a simple strategy that allows clients to subscribe to urls and receive
updates when requests are made by other clients making changes.

To setup spokes do the following in your project:

``` go
import (
  "net/http"
  "os"
	"github.com/nathanborror/gommon/spokes"
)

func main() {
	go spokes.Hub.Run()

  http.Handle("/ws", spokes.SpokeHandler)

  err := http.ListenAndServe(":8080", nil)
	if err != nil {
		os.Exit(1)
	}
}
```

Then use the provided javascript to interact with the WebSocket.


#### Tokens

Tokens is a package for storing push tokens about any devies an auth.User
may have. It currently uses [github.com/anachronistic/apns](github.com/anachronistic/apns)
for sending push notifications to iOS devices.

Just add the http handler like so:

``` go
import (
  "net/http"
  "os"
  "github.com/nathanborror/gommon/tokens"
)

func main() {
  http.Handle("/t/save", tokens.SaveHandler)

  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    os.Exit(1)
  }
}
```

Then in your handler create a list of auth.User.Key and pass it into the
tokenRepo.Push method to send the message to APNS servers. Be sure to put
your push certificate and key in the root directory of your project. You can use
[this tutorial](http://www.raywenderlich.com/32960/apple-push-notification-services-in-ios-6-tutorial-part-1)
to setup your iOS push stuff.

``` go
import (
  "net/http"
  "github.com/nathanborror/gommon/tokens"
)

var tokenRepo = tokens.SqlRepository()
var authRepo = auth.SqlRepository()

func yourHandler() {
  users, err := authRepo.List(100)
  if err != nil {
    panic(err)
  }

  ul := []string{}
  for _, u := range users {
    ul = append(ul, u.Key)
  }

  err = tokenRepo.Push(users, "Your push message", "YourCert.pem", "YourKey.pem")
  if err != nil {
    panic(err)
  }
}
```
