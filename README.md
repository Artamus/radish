# radish
A super barebones Redis clone in Go that implements the following commands:
 - PING
 - ECHO
 - GET
 - SET (with expiry in both seconds and milliseconds)

Inspired by [Build Your Own Redis](https://rohitpaulk.com/articles/redis-0). I followed this guide in what order to build things in and what to use as the test cases. As it stops with the ECHO command, SET, GET and expiry were added based on the knowledge I acquired from the Redis documentation.

The goal of writing this project was to:
 - practice writing Go,
 - play around with things I've learned, but don't have that good hands-on experience with (TCP),
 - try and practice TDD style development.
