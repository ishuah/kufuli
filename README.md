# kufuli

[![Go Report Card](https://goreportcard.com/badge/github.com/ishuah/kufuli)](https://goreportcard.com/report/github.com/ishuah/kufuli)

Kufuli is a centralized locking system for distributed, highly available systems. Written in Go and powered by gRPC.

## Installation
There's currently no build system for this project, you'll have to build from source.

Ensure you have gPRC installed:
`go get -u google.golang.org/grpc`

Clone this repo and install project dependencies:
`git clone git@github.com:ishuah/kufuli.git`
`dep ensure`

You can now build the project:
`go build`

Run the instance:
`./kufuli`

## Writing clients
Clients in Go are supported using the `github.com/ishuah/kufuli/api` package.

A simple example(`examples/client.go`):

```go
package main

import (
	"log"
	"os"

	"github.com/ishuah/kufuli/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		log.Fatal("You need to provide three arguments")
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := api.NewApiClient(conn)

	action := args[0]
	var response *api.Response

	switch action {
	case "lock":
		response, err = c.RequestLock(context.Background(), &api.Request{Resource: args[1], ServiceID: args[2]})
	case "release":
		response, err = c.ReleaseLock(context.Background(), &api.Request{Resource: args[1], ServiceID: args[2]})
	}

	if err != nil {
		log.Fatalf("Error when calling %s: %s", action, err)
	}

        log.Printf("Response from server, success: %t, error: %s", response.Success, response.Error)
}

```

You can run the above example as follows:
`go run client.go lock disk2 sync`

You create a new client instance via `api.NewApiClient`. The client exposes two function, `RequestLock` and `ReleaseLock`. Both functions accept an `api.Request` pointer and return an `api.Response` and an `error`. Further documentation below.

## Support for other languages
gPRC uses [protocol buffers](https://developers.google.com/protocol-buffers/docs/overview) which allows you to generate client and server interfaces from a `.proto` file. A lot of major languages are [supported](https://grpc.io/docs/).

### Example: Generating a python client
Create a new directory, e.g `kufuli-client-py`.
In this new directory, create a new directory called `api`.
Copy the [api/api.proto](https://github.com/ishuah/kufuli/blob/master/api/api.proto) file to `kufuli-client-py/api`.

Install the python gPRC package:
`pip install grpcio-tools`

Generate the python code:
`python -m grpc_tools.protoc -I.  --python_out=. --grpc_python_out=. api/api.proto`

Create a new file, `kufuli-client-py/client.py` and copy the following code into it:

```py
import grpc
import sys
from api import api_pb2, api_pb2_grpc

def run(args):
	channel = grpc.insecure_channel('localhost:7777')
	stub = api_pb2_grpc.ApiStub(channel)
	action = args[0]

	request = api_pb2.Request(resource=args[1], serviceID=args[2])

	try:
		if action == "lock":
			response = stub.RequestLock(request)
		elif action == "release":
			response = stub.ReleaseLock(request)
		else:
			print "unsupported action {}".format(action)
			sys.exit(0)
		print "Response from server: {}".format(response)
	except grpc.RpcError as e:
		print "error when requesting {}: {}".format(action, e)


if __name__ == '__main__':
	args = sys.argv[1:]
	
	if len(args) < 3:
		print "You need to provide three arguments" 
		sys.exit(0)

	run(args)
```

Finally, run the client (make sure the server instance is running):
`python client.py lock disk0 sync` 


# Documentation

<h2>api.proto</h2>
<h3>Request</h3>     
<table>
    <thead>
        <tr><td>Field</td><td>Type</td><td>Label</td><td>Description</td></tr>
    </thead>
    <tbody>
        <tr>
            <td>resource</td>
            <td><a href="#string">string</a></td>
            <td></td>
            <td><p> </p></td>
        </tr>  
        <tr>
            <td>serviceID</td>
            <td><a href="#string">string</a></td>
            <td></td>
            <td><p> </p></td>
        </tr>
    </tbody>
</table>
<h3>Response</h3>
<table>
    <thead>
        <tr><td>Field</td><td>Type</td><td>Label</td><td>Description</td></tr>
    </thead>
    <tbody>
        <tr>
            <td>success</td>
            <td><a href="#bool">bool</a></td>
            <td></td>
            <td><p> </p></td>
        </tr>
        <tr>
            <td>error</td>
            <td><a href="#string">string</a></td>
            <td></td>
            <td><p> </p></td>
        </tr>
    </tbody>
</table>

<h3>Api</h3>
<table>
    <thead>
        <tr>
            <td>Method Name</td>
            <td>Request Type</td><td>Response Type</td>
            <td>Description</td>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>RequestLock</td>
            <td><a href="#api.Request">Request</a></td>
            <td><a href="#api.Response">Response</a></td>
        </tr>
        <tr>
            <td>ReleaseLock</td>
            <td><a href="#api.Request">Request</a></td>
            <td><a href="#api.Response">Response</a></td>
        </tr>
</tbody>
</table>

<h2>Scalar Value Types</h2>
    <table>
      <thead>
        <tr><td>.proto Type</td><td>Notes</td><td>C++ Type</td><td>Java Type</td><td>Python Type</td></tr>
      </thead>
      <tbody>
          <tr id="double">
            <td>double</td>
            <td></td>
            <td>double</td>
            <td>double</td>
            <td>float</td>
          </tr>
          <tr id="float">
            <td>float</td>
            <td></td>
            <td>float</td>
            <td>float</td>
            <td>float</td>
          </tr>
          <tr id="int32">
            <td>int32</td>
            <td>Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead.</td>
            <td>int32</td>
            <td>int</td>
            <td>int</td>
          </tr>
          <tr id="int64">
            <td>int64</td>
            <td>Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead.</td>
            <td>int64</td>
            <td>long</td>
            <td>int/long</td>
          </tr>
          <tr id="uint32">
            <td>uint32</td>
            <td>Uses variable-length encoding.</td>
            <td>uint32</td>
            <td>int</td>
            <td>int/long</td>
          </tr>
          <tr id="uint64">
            <td>uint64</td>
            <td>Uses variable-length encoding.</td>
            <td>uint64</td>
            <td>long</td>
            <td>int/long</td>
          </tr>
          <tr id="sint32">
            <td>sint32</td>
            <td>Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s.</td>
            <td>int32</td>
            <td>int</td>
            <td>int</td>
          </tr>
          <tr id="sint64">
            <td>sint64</td>
            <td>Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s.</td>
            <td>int64</td>
            <td>long</td>
            <td>int/long</td>
          </tr>
          <tr id="fixed32">
            <td>fixed32</td>
            <td>Always four bytes. More efficient than uint32 if values are often greater than 2^28.</td>
            <td>uint32</td>
            <td>int</td>
            <td>int</td>
          </tr>
          <tr id="fixed64">
            <td>fixed64</td>
            <td>Always eight bytes. More efficient than uint64 if values are often greater than 2^56.</td>
            <td>uint64</td>
            <td>long</td>
            <td>int/long</td>
          </tr>
          <tr id="sfixed32">
            <td>sfixed32</td>
            <td>Always four bytes.</td>
            <td>int32</td>
            <td>int</td>
            <td>int</td>
          </tr>
          <tr id="sfixed64">
            <td>sfixed64</td>
            <td>Always eight bytes.</td>
            <td>int64</td>
            <td>long</td>
            <td>int/long</td>
          </tr>
          <tr id="bool">
            <td>bool</td>
            <td></td>
            <td>bool</td>
            <td>boolean</td>
            <td>boolean</td>
          </tr>
          <tr id="string">
            <td>string</td>
            <td>A string must always contain UTF-8 encoded or 7-bit ASCII text.</td>
            <td>string</td>
            <td>String</td>
            <td>str/unicode</td>
          </tr>
          <tr id="bytes">
            <td>bytes</td>
            <td>May contain any arbitrary sequence of bytes.</td>
            <td>string</td>
            <td>ByteString</td>
            <td>str</td>
          </tr>
      </tbody>
    </table>
