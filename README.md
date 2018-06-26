# ParticleMSG

A simple JSON messaging protocol for applications.
Written with plugin managers in mind.

## A sample session

```json
C: {"Type": "_register", "Data": {"Name": "pluginName", "Key": "a sha256-hashed key"}}
S: {"Type": "_registered", "Data": {"Name": "pluginName"}}
OR
S: {"Type": "_alreadyRegistered", "Data": {"Name": "pluginName"}}
OR
S: {"Type": "_invalidKey", "Data": {"Key": "a sha256-hashed key"}}
...
S: {"Type": "_ping", "Data": {}}
C: {"Type": "_pong", "Data": {}}
...
C: {"Type": "_message", "Data": {"To": "pluginName", "Message": {"Type": "aMessage", "Data": {"Some": "data"}}}}
S: {"Type": "_message", "Data": {"From": "pluginName", "Message": {"Type": "aMessage", "Data": {"Some": "data"}}}}
...
C: {"Type": "_quit", "Data": {}}
S: {"Type": "_quit", "Data": {"Reason": "Client Quit"}}
```

On disconnect client generates `{"Type": "_disconnect", "Data": {}}` message.

On disconnect server generates `{"Type": "_disconnect", "Data": {"Who": "pluginName"}}` message.

## Demo

Check `test_client/main.go` and `test_server/main.go` for demos.

```bash
$ cd /path/to/particlemsg_go/
$ bash ./genkey.sh pmsg
$ cd test_server
$ bash ./start.sh
Listening on 0.0.0.0:5050
core: &{_register map[Key:f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b Name:core]}
_server: &{_registered map[Name:core]}
core: &{_message map[Message:map[Data:map[baz:quux foo:bar] Type:foo] To:core]}
core: &{foo map[baz:quux foo:bar]}
_server: &{_message map[From:core Message:map[Type:foo Data:map[baz:quux foo:bar]]]}
```