# OSRM Traffic Updater
The **OSRM Traffic Updater** is designed for pull traffic data from **Traffic Proxy(Telenav)** then dump to OSRM required `traffic.csv`. Refer to [OSRM with Telenav Traffic Design](../docs/design/osrm-with-telenav-traffic.md) and [OSRM Traffic Update](https://github.com/Project-OSRM/osrm-backend/wiki/Traffic) for more details.        
We have implemented both `Python` and `Go` version. Both of them have same function(pull data then dump to csv), but the `Go` implementation is about **23 times** faster than `Python` implementation. So strongly recommended to use `Go` implementation as preference.        
- E.g. `6727490` lines traffic of NA region     
    - `Go` Implementation: about `9 seconds`
    - `Python` Implementation: about `210 seconds`    

## Python Implementation
### Requirements
- `python 2.7`
    - `pip install -U pip`
    - `pip install thrift` (`thrift 0.11.0`)

### Usage
```bash
$ cd traffic_updater/python
$ python osrm_traffic_updater.py -h
usage: osrm_traffic_updater.py [-h] [-p PORT] [-c IP] [-f CSV_FILE]

OSRM traffic updater against traffic proxy.

optional arguments:
  -h, --help            show this help message and exit
  -p PORT, --port PORT  traffic proxy listening port (default: 6666)
  -c IP, --ip IP        traffic proxy ip address (default: 127.0.0.1)
  -f CSV_FILE, --csv_file CSV_FILE
                        OSRM traffic csv file (default: traffic.csv)
```

## Go Implementation
### Requirements
- `go version go1.12.5 linux/amd64`
- `thrift 0.12.0`
    - clone `thrift` from `github.com/apache/thrift`, then checkout branch `0.12.0`
- change `thrift` imports in generated codes `gen-go/proxy` 
    - `git.apache.org/thrift.git/lib/go/thrift` -> `github.com/apache/thrift/lib/go/thrift`


### Usage
```bash
$ cd $GOPATH
$ go install github.com/Telenav/osrm-backend/traffic_updater/go/osrm_traffic_updater
$ ./bin/osrm_traffic_updater -h
Usage of ./bin/osrm_traffic_updater:
  -c string
        traffic proxy ip address (default "127.0.0.1")
  -d    use high precision speeds, i.e. decimal (default true)
  -f string
        OSRM traffic csv file (default "traffic.csv")
  -p int
        traffic proxy listening port (default 6666)
```

