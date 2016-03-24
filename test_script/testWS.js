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
	client.bonbon = function () {
		client.send(JSON.stringify({Cmd: "bonbon"}))
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
	client.history = function(id, number, when) {
		client.send(JSON.stringify({Cmd: "history", With_who: id, Number: number, When: when}))
	}
	client.read = function(id) {
		client.send(JSON.stringify({Cmd: "read", With_who: id}))
	}
	return client
}

const SOCKET_NUM = 5
clients = []

for (var i = 1; i <= SOCKET_NUM; i++) {
	clients[i] = createClient(i)
}

// clients[1].connect("stranger")
// clients[2].connect("stranger")
// clients[1].sendTo(2, "11111111111")
// clients[1].sendTo(2, "22222222222")
// clients[1].sendTo(2, "33333333333")
// clients[1].sendTo(2, "44444444444")
// clients[2].history(2, 3, Math.pow(2, 62))

// rl.on('line', function(line){
// 	id = parseInt(line.split(' ')[0])
// 	msg = line.split(' ').slice(1).join(' ')
// 	clients[id].send(msg)
// })
