# Gommon toolkit

Gommon is a set of packages for quickly developing Go applications.
**This stuff is very new and subject to lots of change. Please fork until
things settle down. Thanks :)**

## Auth

Auth is a simple way to store and authenticate users. It builds on Gorilla
toolkit's sessions package (http://www.gorillatoolkit.org/pkg/sessions) and
uses SQLite to store user information.

Todos:

- Need to save passwords in a more secure way. MD5 isn't ideal.
- Need to write tests

## Render

Render is a collection of basic functions that help you render template
output. If the request is an XHR request it will return a JSON response instead
of rendering it's given template.

There's also a handy filter to help with rendering markdown inside templates
using the Blackfriday package (github.com/russross/blackfriday).
