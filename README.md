# prat
Simple chat server/client in Go using TCP

## Flags
* `--server` <br>
  Start a new server on default address (localhost:9876)
* `--port <portnumber>` <br>
  Set port (default: 9876)
* `--host <hostname>` <br>
  Set host (defautl: localhost)
* `--log <path>` <br>
  Set path where server should output log. If not set, log will output to stdout
* `--silent` <br>
  If silent flag is set, no log output will be generated (either to file or to stdout)

## Example
### Start server
```
prat --server --host localhost --port 9999
```
### Start client
```
prat --host localhost --port 9999
```

## Client commands
* `/setname <name>` <br>
  Sets your username
* `/help`<br>
  Prints the help screen
* `/exit` <br>
  Terminates the client session
