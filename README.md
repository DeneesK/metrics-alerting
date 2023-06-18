# Server for collecting metrics and alerting

Implemented agent and server

### Flags must be specified to run:

#### server's flags:
- `-a` - address and port to run server(default: `localhost:8080`)
- `-l`- logging level(default: `info`)
- `-f` - path to store file containing metrics in JSON(default: `tmp/metrics-db.json`)
- `-r` - load saving data(default: true)
- `-a` - interval of storing data on disk(default: 5 sec)
- `-d` - database's dsn connection configs(default: empty)

#### agent's flags:
- `-a` - address and port to run server(default: `localhost:8080`)
- `-r`- interval of sending metrics to the server(default: 10 sec)
- `-p` - nterval of polling metrics from the runtime(default: 2 sec)

The service allows:
- to collect, store and display metrics
- set up alerts

What is going to be added:
- sending notifications under predefined conditions
