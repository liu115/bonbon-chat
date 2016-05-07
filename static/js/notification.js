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

function NewMessage(who, msg) {
  if (Notification && Notification.permission === "granted" && !isFocus) {
    var n = new Notification('New Message', {
      icon: '',
      body: who + ':' + msg
    });
    n.onshow = function () {
      setTimeout(n.close.bind(n), 5000);
    }
	MSG_SOUND.play();
  }
  return 0;
}
