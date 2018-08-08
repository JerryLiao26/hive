<div align=center>
  <img src='https://raw.githubusercontent.com/JerryLiao26/hive/master/static/hive_small.png' alt='hive logo'>
</div>

# hive

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
Rename to built file to hive, then you shall run server with command
```
hive start 0.0.0.0:12580
```
You shall see a welcome page by visiting ```http://0.0.0.0:12580```. More command can be found with command ```hive help```

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
