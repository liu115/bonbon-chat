// 目前的database放置於/tmp/bonbon.db，可在此建立假資料以供測試，若不懂SQL語法建議使用sqlitestudio有不錯的gui可用

function createSocket(token) {
  var host = window.location.host;
  // 目前使用的 /test/chat/:id 會讓你連進去之後不經確認就把你當做是id為1的使用者
  var chatSocket = new WebSocket("wss://" + host + "/chat/" + token);
  chatSocket.handlers = {}
  // 請見 docs/API0.1 以充分了解 API ，對前端行為不清楚時我們再進行討論，大部份都是顯而易見的
  // addHandler 會可註冊當cmd_type發生時，所要呼叫的函式，請在 chat.js 中註冊相應的行為
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
