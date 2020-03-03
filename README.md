slack-term
==========

A [Slack](https://slack.com) client for your terminal.

![Screenshot](/screenshot.png?raw=true)

Installation
------------

#### Binary installation

[Download](https://github.com/acaloiaro/slack-term/releases) a
compatible binary for your system. For convenience, place `slack-term` in a
directory where you can access it from the command line. Usually this is
`/usr/local/bin`.

```bash
$ mv slack-term /usr/local/bin
```

#### Via Go

If you want, you can also get `slack-term` via Go:

```bash
$ go get github.com/acaloiaro/slack-term
$ cd $GOPATH/src/github.com/acaloiaro/slack-term
$ go install .
```

Setup
-----

1. Get a slack token, click [here](https://api.slack.com/docs/oauth-test-tokens) 

2. Create a `.slack-term` file, and place it in your home directory. Below is
   an example of such a file. You are only required to specify a
   `slack_token`. For more configuration options of the `.slack-term` file,
   see the [wiki](https://github.com/erroneousboat/slack-term/wiki).

```javascript
{
    "slack_token": "yourslacktokenhere"
}
```

Usage
-----

When everything is setup correctly you can run `slack-term` with the following
command: 

```bash
$ slack-term
```

Default Key Mapping
-------------------

Below are the default key-mappings for `slack-term`, you can change them
in your `.slack-term` file.

| mode    | key       | action                       |
|---------|-----------|------------------------------|
| command | `i`       | insert mode                  |
| command | `/`       | search mode                  |
| command | `k`       | move channel cursor up       |
| command | `j`       | move channel cursor down     |
| command | `enter`   | select channel               |
| command | `g`       | move channel cursor top      |
| command | `G`       | move channel cursor bottom   |
| command | `pg-up`   | scroll chat pane up          |
| command | `ctrl-b`  | scroll chat pane up          |
| command | `ctrl-u`  | scroll chat pane up          |
| command | `pg-down` | scroll chat pane down        |
| command | `ctrl-f`  | scroll chat pane down        |
| command | `ctrl-d`  | scroll chat pane down        |
| command | `m`       | toggle message ID visibility |
| command | `n`       | next search match            |
| command | `N`       | previous search match        |
| command | `q`       | quit                         |
| command | `f1`      | help                         |
| insert  | `left`    | move input cursor left       |
| insert  | `right`   | move input cursor right      |
| insert  | `enter`   | send message                 |
| insert  | `esc`     | command mode                 |
| search  | `esc`     | command mode                 |
| search  | `enter`   | command mode                 |

Slash Commands
--------------
| command    | first_param | second param      | description                                                    |
|------------|-------------|-------------------|----------------------------------------------------------------|
| `/delete`  | `msgID`     | `N/A`             | Deletes the message identified by `msgID` from channel history |
| `/edit`    | `msgID`     | `N/A`             | Edit the message identified by `msgID`                         |
| `/thread`  | `msgID`     | `your message`    | Sends threaded message under message identified by `msgID`     |

Note: Use `m` in command mode to toggle message IDs on and off in channels. 

Example Config 
--------------
`~/.slack-term`
```
{
  "slack_token": "<YOUR TOKEN HERE",
  "sidebar_width": 2,
  "notify": "mention",
  "emoji": true,
  "show_unread_only": true,
  "search_timeout": 500,
  "new_message_bell": false,
  "theme": {
    "message": {
      "time_format": "02/01 15:04",
      "time": "fg-green,fg-bold",
      "name": "colorize,fg-bold",
      "text": "fg-blue"
    },
    "channel": {
      "prefix": "fg-red,fg-bold",
      "icon": "fg-green,fg-bold",
      "text": "fg-blue,fg-bold"
    }
  }
}```
