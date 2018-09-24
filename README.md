# Overview
Timing attacks use time as a side-channel which leaks information about the data
that a program is operating on. To illustrate such a timing vulnerability for
educational purposes, this project provides an insecure web server that
does user authentication in non-constant time:

* The server computers a user's unique authentication tag based on a secure
primitive for message authentication codes.
* This tag is then byte-by-byte compared to a user-supplied authentication tag
such that an error is returned as soon as there is a mis-match.

The last step is not constant time, and can be exploited to guess a user's
authentication tag. By default, the server listens for HTTP GET requests on
  `http://localhost:20000/auth/<delay>/<user>/<tag>`,
where `tag` is a 4-byte hex-encoded tag for `<user>`. The `<delay>` specifies
how long the server will pause in ms after each byte-by-byte comparison, making
it easier to exploit the vulnerability without an excessive amount of
repetitions. If access is granted, HTTP status 200 OK is returned.

Access granted example:
  http://localhost:20000/auth/50/alice/c3d36f5f

Access denied example:
  http://localhost:20000/auth/1/alice/c3d36f5f

For further configuration options, e.g., to increase the tag size, invoke the
help flag (`-h` or `--help`).
