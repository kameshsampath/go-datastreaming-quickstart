# Get Started Building Streaming App With Go

A demo to show how to use [go](https://go.dev) and [redpanda](https://redpanda.com) to build simple streaming API application using [Echo Project](https://echo.labstack.com).

## Pre-requisites

- [Docker for Desktop](https://www.docker.com/products/docker-desktop/)
- [rpk CLI](https://docs.redpanda.com/current/get-started/rpk-install/)
- [httpie](https://httpie.io)

## Run Application

```shell
go run server.go -brokers "$RPK_BROKERS"
```

### Start Redpanda Server

```shell
rpk container start -n 1
```

**NOTE**:

> Make a note of the RPK_BROKERS value and export as instructed

### Create `greetings` Topic

```shell
rpk topic create greetings --partitions=3 --replicas=1
```

### Start Streaming Consumer

Open a new terminal window

```shell
http --stream ':8080/?topic=greetings'
```

Since the application is set to use the same Consumer Group ID of format `<topic>-echo-group` the consumer will not read the old messages when restarted. If you wish to read all from start run the following command,

```shell
http --stream ':8080/?topic=greetings&group=new'
```

This command will create a new consumer group on each request.

### Produce a message

Open new terminal and run the following command to produce a new message to topic called `greetings`,

```shell
http -v :8080/ topic=greetings key='1' value='Hello World!'
```

## Clean up

```shell
rpk container purge
```
