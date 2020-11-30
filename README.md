# emsk
A simple cli for interacting with MSK and Lambda

Lambda consumers of MSK topics have UUID's for consumer group names, which make it difficult to identify which lambda consumers are configured to consume messages from a certain topic


### Requirements
[Install](https://formulae.brew.sh/formula/go) `Go 1.15`

### Building binary

Clone repo

```
git clone git@github.com:colincgn/emsk.git
```

from `emsk` folder.

```
go build
```

This will create an executable binary with a few handy commands when working with MSK.

#### Configuration

For convenience, you can use some environment variables instead of passing in command line flags on each call.

```
EMSK_BOOTSTRAP_SERVERS=1-aws-msk-bootstrap-server:9094,2-aws-msk-bootstrap-server:9094
EMSK_TLS_ENABLED=TRUE
```


#### List all topics
```
./emsk topic list
```

#### List all consumer groups
```
./emsk consumergroup list
``` 

example output.

```
{
  "Id": "23c2bf3d-91ea-44bc-8ab1-298f758d4fa7",
  "ActiveMembers": 1,
  "Members": [
    {
      "ClientId": "23c2bf3d-91ea-44bc-8ab1-298f758d4fa7",
      "Topics": [
        "example-topic-one"
      ]
    }
  ],
  "LastKnownStatus": "OK",
  "FunctionArn": "arn:function:example-topic-one-dev-ExampleConsumer"
}
```
