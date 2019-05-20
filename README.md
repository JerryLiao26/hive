<div align=center>
  <img src='https://raw.githubusercontent.com/JerryLiao26/hive/master/static/hive_small.png' alt='hive logo'>
</div>

# hive

[![Build Status](https://travis-ci.org/JerryLiao26/hive.svg?branch=master)](https://travis-ci.org/JerryLiao26/hive)
[![Go Report Card](https://goreportcard.com/badge/github.com/JerryLiao26/hive)](https://goreportcard.com/report/github.com/JerryLiao26/hive)
[![codebeat badge](https://codebeat.co/badges/4e16d63d-1c79-43ed-86af-ac6fdb7c0367)](https://codebeat.co/projects/github-com-jerryliao26-hive-master)
[![License](https://img.shields.io/github/license/JerryLiao26/hive.svg)](https://opensource.org/licenses/MIT)

**Hive** is the light-weight and elegant message center for developers

## Origin

Originally [LAM](https://github.com/JerryLiao26/LAM), **hive** is now a more elegant way to collect your messages.

## Use

**Hive** is under heavy develop and has no official release version at the present. If you would like to try it now, try [Develop](#develop) option shown below.

### Pushers
- [hive-pusher-js](https://github.com/JerryLiao26/hive-pusher-js)
- [hive-pusher-go](https://github.com/JerryLiao26/hive-pusher-go)

## Develop

### Prepare

**Hive** use go [mod](https://github.com/golang/go/wiki/Modules) for dependency management. Run ```go build``` in project folder will install all dependencies under ```$GOPATH/pkg/mod/``` folder.

### Database

To use **hive**, you need a database named ```hive```, and tables as follows:
- message
  - id(INT) [PRIMARY, AUTO_INCREMENT]
  - tag(VARCHAR)
  - admin(VARCHAR)
  - content(TEXT)
  - timestamp(DATETIME)

- token
  - tag(VARCHAR) [PRIMARY]
  - token(VARCHAR)
  - admin(VARCHAR)
  - timestamp(DATETIME)

- admin
  - name(VARCHAR) [PRIMARY]
  - token(VARCHAR)
  - timestamp(DATETIME)

There's an sql file under ```sql/``` for you to import

### Build and test

Navigate to project folder, simply run
```
go build main.go
```
Rename the built file to hive, then you shall run server with command
```
hive start 0.0.0.0:12580
```
You shall see a welcome page by visiting ```http://0.0.0.0:12580/hello```. More commands can be found with command ```hive help```

## Pages

### /hello
Welcome and navigation page of **hive**

### /auth
Enter your admin token to grant you permission for using ```/dashboard``` and ```/panel```, session is valid for 24 hours

### /panel
Show server status, informations about this **hive** and so on

### /dashboard
Show tags of the authorized admin, as well as messages. Dashboard automatically refresh message list every 5 seconds

## APIs

### /auth
- Method: POST
- Require:
```json
{
  "token": "Your token here"
}
```
- Respond:
```json
{
  "code": 200,
  "text": "Token verified",
  "method": "POST"
}
```

### /info
- Method: POST
- Require:
```json
{
  "sessionId": "Your session ID here"
}
```
- Respond:
```json
{
  "code": 200,
  "text": "Session verified",
  "adminName": "Your admin name",
  "startTime": "2018-08-08 18:08:08",
  "serverOS": "linux",
  "serverArch": "amd64",
  "method": "POST"
}
```

### /session
- Method: POST
- Require:
```json
{
  "sessionId": "Your session ID here"
}
```
- Respond:
```json
{
  "code": 200,
  "text": "Session verified",
  "method": "POST"
}
```

### /messages
- Method: POST
- Require:
```json
{
  "sessionId": "Your session ID here"
}
```
- Respond:
```json
{
  "code": 200,
  "text": "Session verified",
  "messages": [
    {"id": 0, "tag": "wow", "admin": "test", "content": "hello", "timestamp": "2018-08-09T16:22:59Z"}
  ],
  "method": "POST"
}
```

### /send
- Method: POST
- Require:
```json
{
  "text": "Your message here",
  "token": "Your token here"
}
```
- Respond:
```json
{
  "code": 200,
  "text": "Message from Your tag saved",
  "method": "POST"
}
```

## To-dos

- [ ] Code pushers and command-line pusher
- [ ] Add super admin and a button to restart **hive** for him/her
- [ ] Add limit for API request(e.g. 60/min)
- [ ] Improve token generate methods
- [ ] Running with goroutine
- [ ] Server daemon
