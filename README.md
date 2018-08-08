<div align=center>
  <img src='https://raw.githubusercontent.com/JerryLiao26/hive/master/static/hive_small.png' alt='hive logo'>
</div>

# hive

[![Build Status](https://travis-ci.org/JerryLiao26/hive.svg?branch=master)](https://travis-ci.org/JerryLiao26/hive)
[![Go Report Card](https://goreportcard.com/badge/github.com/JerryLiao26/hive)](https://goreportcard.com/report/github.com/JerryLiao26/hive)
[![codebeat badge](https://codebeat.co/badges/3737204a-ce9b-4d80-a930-5a8735fd709c)](https://codebeat.co/projects/github-com-jerryliao26-hive-master)
[![License](https://img.shields.io/github/license/JerryLiao26/hive.svg)](https://opensource.org/licenses/MIT)

hive is the light-weight and elegant message center for developers

## Origin

Originally [LAM](https://github.com/JerryLiao26/LAM), hive is now a more elegant way to collect your messages.

## Use

hive is under heavy develop and has no official release version at the present. If you would like to try it now, try [Develop](#develop) option shown below.

## Develop

### Prepare

hive use [dep](https://golang.github.io/dep/) as Go package manager. Run ```dep ensure``` in project folder will install all dependencies under ```vendor/``` folder.

### Database

To use hive, you need a database named ```hive```, and tables as follows:
- message
  - id(INT) [PRIMARY, AUTO_INCREMENT]
  - tag(VARCHAR)
  - content(TEXT)
  - timestamp(TIMESTAMP)
  - ifRead(TINYINT)

- token
  - tag(VARCHAR) [PRIMARY]
  - token(VARCHAR)
  - timestamp(TIMESTAMP)

There's an sql file under ```sql/``` for you to import

### Build and test

Navigate to project folder, simply run
```
go build src/*
```
Rename the built file to hive, then you shall run server with command
```
hive start 0.0.0.0:12580
```
You shall see a welcome page by visiting ```http://0.0.0.0:12580/hello```. More command can be found with command ```hive help```

## APIs

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
