Power start
===========

Server for starting and stopping machines in the LAN.

Server sent special magic packet to launch computer via Wake-On-LAN.

For each machine have stored count of start requests. If start requests is 0 then server will sent shutdown request to the machine. But for some time it will keep machine on when start requests count will 0.

## API

### Add machine
`curl -X POST -d '{"name":"the_name", "mac":"01:23:45:67:89:AB"}' localhost:3000/add`

### Delete machine
`curl -X POST localhost:3000/remove/the_name`

### Start request
`curl -X POST localhost:3000/start/the_name`

### Stop reqest
`curl -X POST localhost:3000/stop/the_name`

### List of machines
`curl localhost:3000/list`

Example of response
```json
[
  {
    "name": "the_name_1",
    "mac": "01:23:45:67:89:AB",
    "requests": 0,
    "isRunning": false
  },
  {
    "name": "the_name_2",
    "mac": "FE:DC:BA:98:76:54",
    "requests": 2,
    "isRunning": true
  }
]
```