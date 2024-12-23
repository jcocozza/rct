# rct

rct stands for remote copy text.

I've found that when I'm working in a remote environment, there hasn't been a good way to copy text from remote to local.
This project is an attempt to solve this problem.

The requirements for this project is as follows:

1. Copy from remote to local should be simple. `rct <my text>`
2. Exclusive remote/client relationship (only verified machines can send data)
3. End-to-End encryption

## Setup

rct will look for an `.rst.json` config file in your home directory.

Here is a sample config.

The `server` component is where rct will listen for incoming tcp connections.
Specify the token to limit incoming tcp data that includes the correct token.

The `delivery` component is where rct will send data.
This can be more then one address which allows you to send data to more then one machine.
(how practical this is, is up for debate, but nevertheless, it is an options)

```json
{
  "server": {
    "addr": "10.0.0.6:8080", # this is the address of your local machine and the port to listen on
    "token": ""
  },
  "delivery": [ # these addresses are for sending data
    {
      "addr": "10.0.0.6:8080",
      "token": ""
    }
  ]
}
```

Because rct is written for both remote and local, both `server` and `delivery` are in the config file.
However, you will really only care about either `server` or `delivery` depending on the machine you are on.
With the intended rct use case:
- locally: you will care about the `server`
- remote: you will care about `delivery`

## Architecture

rct communicates over tcp.
The cli is equipped to act as both a server and a client.

The client pushes data to the server where it is copied to the system clipboard.

As such, the intended use of rct is running the rct server on your local machine as a background process with the rct client on the remote.
(Typically, its very easy to copy locally and paste into remote. The reverse is usually not so easy.)

```
+-------------------+                   +------------------+
|       Local       | <---- (tcp) ----> |      remote      |
|  (rct tcp server) |                   | (rct tcp client) |
+-------------------+                   +------------------+
```
