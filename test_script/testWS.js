var WS = require('ws')
var util = require('util')
var readline = require('readline')
var rl = readline.createInterface({
	input: process.stdin,
	output: process.stdout,
	terminal: false
});

function createClient(id) {
	url = "ws://localhost:8080/test/chat/" + id.toString()
	var client = new WS(url);
	client.id = id
	client.receive = []
	client.on('message', function(msg) {
		client.receive.push(msg)
	})

	client.on('open', function () {
	})

	client.connect = function (type) {
		client.send(JSON.stringify({Cmd: "connect", Type: type}))
	}
	client.disconnect = function () {
		client.send(JSON.stringify({Cmd: "disconnect"}))
	}
	client.sendTo = function (id, msg) {
		client.send(JSON.stringify({Cmd: "send", Who: id, Msg: msg}))
	}
	client.toStranger = function (msg) {
		client.sendTo(0, msg)
	}
	return client
}

const SOCKET_NUM = 5
clients = []

for (var i = 1; i <= SOCKET_NUM; i++) {
	clients[i] = createClient(i)
}

// rl.on('line', function(line){
// 	id = parseInt(line.split(' ')[0])
// 	msg = line.split(' ').slice(1).join(' ')
// 	clients[id].send(msg)
// })
