# rct

rct stands for remote copy text.

I've found that when I'm working in a remote environment, there hasn't been a good way to copy text from remote to local.
This project is an attempt to solve this problem.
(In particular, I often encounter this issue when using vim on remote servers. [rct.vim](https://github.com/jcocozza/rct.vim) is a vim plugin that wraps rct to solve this.)

The requirements for this project is as follows:

1. Copy from remote to local should be simple.

- on remote, simply call `rct <my text>` (provided `rct listen` has been called on local)

2. Exclusive remote/client relationship (only verified machines can send data)

- the config file supports adding a token. This essentially acts as an API key that needs to be included with any request to the rct server.

## Setup

rct will look for an `.rst.json` config file in your home directory.

A config can be created with `gen-config` command.

Here is a sample config.

The `server` component is where rct will listen for incoming tcp connections.
Specify the token to limit incoming tcp data that includes the correct token.

The `delivery` component is where rct will send data.
This can be more then one address which allows you to send data to more then one machine.
(the practicality of this is up for debate, but nevertheless it is an option)

```json
{
  "server": {
    "addr": "10.0.0.6:8080", # the address that rct will run the listen server (in this case, 10.0.0.6 is the address of the local machine)
    "token": ""
  },
  "delivery": [ # these addresses are for sending data
    {
      "addr": "10.0.0.8:8080",
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

I recommend adding rct to your path. (For example, by moving it to `/usr/local/bin`)
From there, so you don't need to keep restarting the listener you can configure your local machine to run `rct listen -d` on start up.

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
|                   |        |    (rct server)   |                   |   (rct client)   |
|     (clipboard)   |        |     (10.0.0.6)    |                   |    (10.0.0.8)    |
+-------------------+        +-------------------+                   +------------------+
```

rct communicates over tcp.
The cli is equipped to act as both a server and a client.

The client pushes data to the server where it is copied to the system clipboard.

As such, the intended use of rct is running the rct server on your local machine as a background process with the rct client on the remote.
(Typically, it is very easy to copy locally and paste into remote. The reverse is usually not so easy.)

Where possible, rct tries to work with posix tools.
For example `cat my_file.txt | rct` works just fine as a way to copy the text of a file to the local machine.

## Example

Using the diagram above, we have 2 config files (one on each machine).
Notice that delivery on one matches server on the other.

1. On remote (10.0.0.8)

```json
{
  "server": {
    "addr": "",
    "token": ""
  },
  "delivery": [
    {
      "addr": "10.0.0.6:8080",
      "token": ""
    }
  ]
}
```

2. On local (10.0.0.6)

```json
{
  "server": {
    "addr": "10.0.0.6:8080",
    "token": ""
  },
  "delivery": [
    {
      "addr": "",
      "token": ""
    }
  ]
}
```
