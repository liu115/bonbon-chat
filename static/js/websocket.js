function createSocket() {
  var host = window.location.host;
  var chatSocket = new WebSocket("ws://" + host + "/test/chat/1");
  chatSocket.handlers = {}
  chatSocket.addHandler = function (cmd_type, handler) {
    if (chatSocket.handlers[cmd_type] == null) {
      chatSocket.handlers[cmd_type] = [handler]
    } else {
      chatSocket.handlers[cmd_type].push(handler)
    }
  }
  chatSocket.onmessage = function(e) {
    cmd = JSON.parse(e.data)
    var hs = chatSocket.handlers[cmd.Cmd]
    if (hs == null) {
      return
    }
    for (var i = 0; i < hs.length; i++) {
      hs[i](cmd)
    }
    console.log(e.data);
  }
  return chatSocket
}
