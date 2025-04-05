## TODO

* Implement websockets using goolang

* on_message

## live

* users current pos will be a highlight in the server/client

* implement the admin_code to be the connection to the admin, each time a new client connects to the server, a new requeste to the admin is made, returning his buffer lines

* connect: client,
    * on_connect create floatwindow, send message to server (client x connected), then receive response with buffer lines and file type
    * on_message send message to server (bufferlines)

* startServer: server
    * on_start send bufnum
