# Gommon toolkit

Gommon is a set of packages for quickly prototyping Go applications.
**This stuff is very new and subject to lots of iteration. Feel free to look
and riff but don't expect any sort of stability.**

## Auth

Auth is a simple and dumb way to store and authenticate users. It builds on the
Gorilla toolkit's sessions package (http://www.gorillatoolkit.org/pkg/sessions)
and uses SQLite to store user information. Please note, passwords are currently
stored in a very insecure fashion.

## Render

Render is a collection of basic functions that help render template output.
If the request is an XHR request it will return a JSON response instead
of rendering it's given template. You can also append `?json` to any URL to
force a JSON response.

## Spokes

Spokes is a basic WebSocket pub/sub implemetation using the Gorilla WebSocket
toolkit. It's based off the [chat example](https://github.com/gorilla/websocket/tree/master/examples/chat)
but adds a simple strategy that allows clients to subscribe to urls and receive
updates when requests are made by other clients making changes.

## Tokens

Tokens is a package for storing push tokens about any devies an auth.User
may have. It currently uses [github.com/anachronistic/apns](github.com/anachronistic/apns)
for sending push notifications to iOS devices.
