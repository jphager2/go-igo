# igo

The game of go implemented in Golang.

## On Scoring (TODO)

[Static Scoring Algorithm][0]

## On Input (TODO)

Currently this program runs on the command line interactively. Ideally it would
run as a daemon and accept [GTP (Go Text Protocol)][1] on a known socket.

## On Web (TODO)

To run this program for the web, it would probably need to be rewritten in a
way that state can be managed on the server side in a persistant way. The
single daemon described above would be inefficient in a web context and also
would make multiple concurrent games complicated. If websockets would be used
for communication after the initial connection, then probably the game could be
running in a go routine that would hold the websocket connection.

[0]: https://www.oipaz.net/Carta.pdf
[1]: https://www.lysator.liu.se/~gunnar/gtp/gtp2-spec-draft2/gtp2-spec.html
