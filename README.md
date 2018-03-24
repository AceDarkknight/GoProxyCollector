# GoProxyCollector
## Introduction
GoProxyCollector is a proxy collector written in Go.
## Advantage
- Out-of-box: Use [boltdb](https://github.com/boltdb/bolt) as embedded storage. No need to install extra database.
- Configurable: Easy to support more proxy website by add your own rule in configuration.
## Installation
- Download
```
go get -u github.com/AceDarkkinght/GoProxyCollector
```
- Start up

```
go build GoProxyPool
```
## Usage
- Get a ip

```
GET http://localhost:8090/get
```
The respone is json.
```json
{
   "ip":"1.180.235.165",
   "port":3128,
   "location":"内蒙古 电信",
   "source":"http://www.ip181.com/",
   "speed":0.13
}
```
- Delete
```
GET http://localhost:8090/delete?ip=1.2.3.4
```
