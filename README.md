# rct

rct stands for remote copy text.

I've found that when I'm working in a remote environment, there hasn't been a good way to copy text from remote to local.
This project is an attempt to solve this problem.
(In particular, I often encounter this issue when using vim on remote servers. [rct.vim](https://github.com/jcocozza/rct.vim) is a vim plugin that wraps rct to solve this.)

The requirements for this project is as follows:

1. Copy from remote to local should be simple.

- on remote, simply call `rct <my text>` (provided `rct listen` has been called on local)

2. Exclusive remote/client relationship (only verified machines can send data)

- the config file supports including a token which all data will have to be sent with

## Setup

rct will look for an `.rst.json` config file in your home directory.

Here is a sample config.

The `server` component is where rct will listen for incoming tcp connections.
Specify the token to limit incoming tcp data that includes the correct token.

The `delivery` component is where rct will send data.
This can be more then one address which allows you to send data to more then one machine.
(the practicality of this is up for debate, but nevertheless it is an option)

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

- locally: you will care about `server`
- remote: you will care about `delivery`

## Basic Flow

Assuming everything is set up properly, basic use is quite simple:

1. On local machine:

```bash
rct listen -d # run the rct listener in the background
```

2. Working on remote:

```bash
rct 'my message' # this will deliver 'my message' to the local machine's clipboard
```

## Design

```
+-------------------+        +-------------------+                   +------------------+
|       local       |  <---  |       local       | <---- (tcp) ----> |      remote      |
|     (clipboard)   |        |    (rct server)   |                   |   (rct client)   |
+-------------------+        +-------------------+                   +------------------+
```

rct communicates over tcp.
The cli is equipped to act as both a server and a client.

The client pushes data to the server where it is copied to the system clipboard.

As such, the intended use of rct is running the rct server on your local machine as a background process with the rct client on the remote.
(Typically, it is very easy to copy locally and paste into remote. The reverse is usually not so easy.)
