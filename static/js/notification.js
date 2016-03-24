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

function NewMessage(who, msg) {
  if (Notification && Notification.permission === "granted") {
    var n = new Notification('New Message', {
      icon: '',
      body: who + ':' + msg
    });
    n.onshow = function () {
      setTimeout(n.close.bind(n), 5000);
    }
  }
  return 0;
}
