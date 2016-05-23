window.addEventListener('load', function () {
  //Ask for permission
  if (Notification && Notification.permission !== "granted") {
    Notification.requestPermission(function (status) {
      if (Notification.permission !== status) {
        Notification.permission = status;
      }
    });
  }
});

var isFocus = true;

window.onfocus = function () {
  isFocus = true;
}

window.onblur = function () {
  isFocus = false;
}

var MSG_SOUND = new Audio("/static/audio/msg.wav");
var BON_SOUND = new Audio("/static/audio/bonbon.mp3");

function NewMessage(who, msg) {
  if (Notification && Notification.permission === "granted" && !isFocus) {
    var n = new Notification(who, {
      icon: '',
      body: msg
    });
    n.onshow = function () {
      setTimeout(n.close.bind(n), 5000);
    }
	MSG_SOUND.play();
  }
  return 0;
}

function NewFriend() {
  if (Notification && Notification.permission === "granted" && !isFocus) {
    var n = new Notification('成為朋友', {
      icon: '',
      body: '已經與匿名對象成為好友囉'
    });
    n.onshow = function () {
      setTimeout(n.close.bind(n), 5000);
    }
	BON_SOUND.play();
  }
}
