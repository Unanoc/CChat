# clichat
Command Line Interface chat which is based on sockets

## Install
```
git clone https://github.com/leedaniil/clichat.git  (https)
```
or
```
git clone git@github.com:leedaniil/clichat.git  (ssh)
```

## How to run the server?
```
cd clichat/cmd/chat-server && go run main.go
```

## How to run the client?
```
cd clichat/cmd/chat-client && go run main.go
```

## Manual for client
First of all, the client should enter the name and the name of the room to which he wants to connect.
If the room does not exists, server will create room with this name.
If the name of client is not uniq, server will ask to enter another name.

After successful connection the client recieves history of room's messages (last 128 messages).

Client is able to write and read messages. 
Besides, there are 3 commands for client:

* **/quit** - exit
* **/list** - get list of clients in room
* **/change_room** - change room


## How it looks like?
<img src="https://github.com/leedaniil/clichat/blob/master/src/clients.png" width="720">
