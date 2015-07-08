var chatSocket;
var msgBox = document.getElementById("message-content");

function createSocket() {
    _chatSocket = new WebSocket("ws://mros.csie.org:8080/chatRoomSocket", "my-custom-protocol");

    _chatSocket.onopen = function() {
        msgBox.innerHTML += "<p><span class='message-from-system'>已建立連線</span></p>";
    }

    _chatSocket.onmessage = function(msg) {
        var jmsg = JSON.parse(msg.data);
        if(jmsg["me"]) {
            console.log("me");
            msgBox.innerHTML += "<p><span class='message-from-me'>" + jmsg["msg"] + "</span></p>";
        }
        else {
            console.log("others");
            msgBox.innerHTML += "<p><span class='message-from-others'>" + jmsg["msg"] + "</span></p>";
        }
    };

    _chatSocket.onclose = function() {
        msgBox.innerHTML += "<p><span class='message-from-system'>已中斷連線</span></p>";
    }
    return _chatSocket;
}

chatSocket = createSocket();

var btnNew = document.getElementById("new-connection");
var btnSend = document.getElementById("button-send-message");
var msgSend = document.getElementById("text-send-message");

btnNew.onclick = function() {
    chatSocket.close();
    chatSocket = createSocket();
    return false;
}

btnSend.onclick = function() {
    chatSocket.send(msgSend.value);
}

msgSend.onkeypress = function(event) {
    if(event.keyCode == 13) {
        chatSocket.send(msgSend.value);
        this.value = "";
    }
}