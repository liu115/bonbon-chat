var React = require('../bower/react/react-with-addons.js');
var MessageBalloon  = require('./message.jsx');
var ReactCSSTransitionGroup = React.addons.CSSTransitionGroup;

var SignClass = React.createClass({
  getInitialState: function() {
    return {
      setting: false,
      value: ''
    };
  },

  handleClick: function() {
    this.setState({setting: true});
    //React.findDOMNode(this.refs.refInput).focus();
  },

  handleType: function(e) {
    var keyInput = e.keyCode == 0 ? e.which : e.keyCode;
    if (keyInput == 13) {
      console.log('trying to set ' + this.state.value + ' as Sign.');
      this.props.chatSocket.send(JSON.stringify({Cmd: "setting", Setting: {Sign: this.state.value}}));
      this.setState({
        setting: false,
        value: ''
      });
    }
  },

  handleChange: function(e) {
    this.setState({
      value: e.target.value
    });
  },

  render: function() {
    if (this.state.setting == true) {
      return (
        <div id="sign-input-wrapper">
          <input type="text" id="sign-input" ref="refInput" value={this.state.value} onKeyPress={this.handleType} onChange={this.handleChange} placeholder="按Enter確認更改簽名"/>
        </div>
      );
    }
    else {
      return (
        <div>
          <a id="profile-status" onClick={this.handleClick}>
            {this.props.sign}
            <i style={{margin: '5px'}} className="fa fa-pencil fa-lg"></i>
          </a>
        </div>
      );
    }
  }
});

var SideBar = React.createClass({
  getInitialState: function() {
    this.props.chatSocket.addHandler("init", function(cmd) {
      this.setState({Sign: cmd.Setting.Sign});
    }.bind(this));
    this.props.chatSocket.addHandler("setting", function(cmd) {
      if (cmd.OK == true) {
        console.log('setting success, new sign is ' + cmd.Setting.Sign);
        this.setState({Sign: cmd.Setting.Sign});
      }
      else {
        console.log('setting failed!');
      }
    }.bind(this));
    return {Sign: "", buttonList: null};
  },
  NewConnection: function(i) {
      console.log(i);
      switch (i) {
        case 1:
          this.props.chatSocket.send(JSON.stringify({Cmd: "connect", Type: "L1_FB_friend"}));
          break;
        case 2:
          this.props.chatSocket.send(JSON.stringify({Cmd: "connect", Type: "L2_FB_friend"}));
          break;
        case 3:
          this.props.chatSocket.send(JSON.stringify({Cmd: "connect", Type: "stranger"}));
          break;
      }
      this.setState({buttonList: null});
  },
  logout: function() {
    localStorage.setItem('login', 'false');
    this.props.logout();
  },
  handleClick: function() {
    if (!this.state.buttonList) {
      this.setState({
        buttonList: (
          <div className="connection-menu">
            <ul id="connection-list">
              <li className="connection-line"><a className="connection-button" onClick={this.NewConnection.bind(this,1)}>FB的好友</a></li>
              <li className="connection-line"><a className="connection-button" onClick={this.NewConnection.bind(this,2)}>朋友的朋友</a></li>
              <li className="connection-line"><a className="connection-button" onClick={this.NewConnection.bind(this,3)}>陌生人</a></li>
              <li className="connection-line"><a className="connection-button" onClick={this.NewConnection.bind(this,4)}>取消</a></li>
            </ul>
          </div>
        )
      });
    }
    else {
      this.setState({buttonList: null});
    }
  },
  render: function() {
    return (
      //<!-- start of navigation area -->
      <nav id="sidebar-panel">
        <div id="sidebar-profile">
          <span id="profile-avatar"><a><img src="/static/img/me_finn.jpg"/></a></span>
          <SignClass sign={this.state.Sign} chatSocket={this.props.chatSocket}/>
        </div>
        <a id="new-connection" onClick={/*this.props.changeState.bind(null, 'new_connection')*/this.handleClick }>建立新連線</a>
        <ReactCSSTransitionGroup transitionName="list-animate">
          {this.state.buttonList}
        </ReactCSSTransitionGroup>
        <ul id="menu">
          <li><a onClick={this.logout}><span><i className="fa fa-sign-out"></i><span style={{margin: '0px'}}>登出</span></span></a></li>
        </ul>
      </nav>
      //<!-- end of navigation area -->
    );
  }
});

var FriendBox = React.createClass({
  handleClick: function() {
    this.props.select(this.props.index);
    this.props.changeState('chat');
  },

  render: function() {
    return (
      <div className={"friend-unit " + "friend-" + this.props.friend.stat + (this.props.friend.online ? '' : " off-line")} onClick={this.handleClick}>
        <div className={(this.props.index == 0) ? "stranger-avatar": "friend-avatar"}>
          <img src={this.props.friend.img}/>
        </div>
        <div className="friend-info">
          <p className="friend-info-name">{this.props.friend.name}</p>
          <p className="friend-info-status">
            {function() {
              if (this.props.friend.messages.length > 0)
              {
                if (this.props.friend.messages[this.props.friend.messages.length - 1].content.length > 20) {
                  return this.props.friend.messages[this.props.friend.messages.length - 1].content.slice(0,20) + '...';
                }
                return (this.props.friend.messages[this.props.friend.messages.length - 1].content);
              }
              else return ('');
            }.bind(this)()}
          </p>
        </div>
        <div style={{clear: "both"}}></div>
      </div>
    );
  }
});

var FriendList = React.createClass({
  getInitialState: function() {
    this.props.chatSocket.addHandler('status', function(cmd) {
    }.bind(this));
    return {
      filterText: ''
    };
  },
  handleFilterInput: function(e) {
    this.setState({
      filterText: e.target.value
    });
  },
  render: function() {
    var friendBoxs = [];
    for (var i = 0; i < this.props.friends.length; i++) {
      if (this.props.friends[i].name.indexOf(this.state.filterText) === -1) continue;
      friendBoxs.push(<FriendBox index={i} friend={this.props.friends[i]} changeState={this.props.changeState} select={this.props.select}/>);
    }
    return (
      <div id="friend-area">
        <div id="friend-search">
          <div id="wrapper-input-search" className="wrapper-input">
            <input type="text" placeholder="搜尋朋友" ref="filterText" onChange={this.handleFilterInput}/>
          </div>
        </div>
        { friendBoxs }
      </div>
    );
  }
});

var ChatRoom = React.createClass({
  getInitialState: function() {
    //name is this.props.name and header take from the name
    this.props.chatSocket.addHandler('history', function(cmd) {
      var node = React.findDOMNode(this.refs.refContent);
      node.scrollTop = this.scrollTop + (node.scrollHeight - this.scrollHeight);
    }.bind(this));
    /* lock when message are sending */
    this.sendLock = 'unlock';
    this.props.chatSocket.addHandler('send', function(cmd) {
      this.sendLock = 'unlock';
      if (cmd.OK == true) {
        this.setState({
          userInput: ''
        });
      }
    }.bind(this));
    return {
      userInput: ''
    };
  },
  handleChange: function(e) {
    this.setState({
      userInput: e.target.value
    });
  },
  componentWillUpdate: function() {
    var node = React.findDOMNode(this.refs.refContent);
    this.scrollHeight = node.scrollHeight;
    this.scrollTop = node.scrollTop;
    this.shouldScrollBottom = node.scrollTop + node.offsetHeight === node.scrollHeight;
  },
  componentDidUpdate: function() {
    if (this.shouldScrollBottom) {
      var node = React.findDOMNode(this.refs.refContent);
      node.scrollTop = node.scrollHeight;
      var messageLength = this.props.friends[this.props.target].messages.length;
      if (messageLength <= 0) return 0;
      var lastMessage = false;
      for (var i = messageLength - 1; i >= 0; i--) {
        if (this.props.friends[this.props.target].messages[i] == 'others')
          lastMessage = this.props.friends[this.props.target].messages[i];
      }
      if (lastMessage === false) return 0;
      if (this.props.friends[this.props.target].read < lastMessage.time) {
        console.log('read auto');
        this.props.chatSocket.send(JSON.stringify({Cmd: "read", With_who: this.props.friends[this.props.target].ID, Time: lastMessage.time}));
      }
    }
  },
  sendMessage: function(e) {
    //send it to websocket
    //this.state.messages.splice(0, 0, ['me', 'lalala']);
    if (this.state.userInput != '' && this.sendLock == 'unlock'){
      this.props.addMessage(this.props.target, 'buttom', {
        from: 'me',
        content: this.state.userInput
      });
      this.sendLock = 'lock';
    }
    this.focusInput();
  },
  sendMessageByKeyboard: function(e) {
    var keyInput = e.keyCode == 0 ? e.which : e.keyCode;
    if (keyInput == 13 && !e.shiftKey) {
      this.sendMessage();
      e.preventDefault();
    }
  },
  focusInput: function() {
    React.findDOMNode(this.refs.refInput).focus();

  },
  componentDidMount: function() {
    this.historyLock = 'unlock'
    React.findDOMNode(this.refs.refInput).focus();
  },
  componentWillUnmount: function() {
  },
  bonbon: function() {
    this.props.chatSocket.send(JSON.stringify({Cmd: "bonbon"}));
  },
  disconnect: function() {
    this.props.chatSocket.send(JSON.stringify({Cmd: "disconnect"}));
  },
  handleScroll: function() {
    var node = React.findDOMNode(this.refs.refContent);
    if (node.scrollTop <= 0) {
      if (this.historyLock === 'unlock') {
        if (this.props.target == 0) return false;
        var who = this.props.friends[this.props.target].ID;
        var time = this.props.friends[this.props.target].messages[0].time;
        console.log(time);
        this.props.chatSocket.send(JSON.stringify({Cmd: "history", With_who: who, Number: 15, When: time, Order: 0}));
      }
    }
    if (node.scrollTop + node.offsetHeight === node.scrollHeight) {
      var messageLength = this.props.friends[this.props.target].messages.length;
      if (messageLength <= 0) return 0;
      var lastMessage = false;
      for (var i = messageLength - 1; i >= 0; i--) {
        if (this.props.friends[this.props.target].messages[i] == 'other')
          lastMessage = this.props.friends[this.props.target].messages[i];
      }
      if (lastMessage === false) return 0;
      if (this.props.friends[this.props.target].read < lastMessage.time) {
        console.log('read');
        this.props.chatSocket.send(JSON.stringify({Cmd: "read", With_who: this.props.friends[this.props.target].ID, Time: lastMessage.time}));
      }
    }
  },
  render: function() {
    return (
      <div id="message-area">
        <div id="message-header" className="area-header"r ref="header">
          {this.props.friends[this.props.target].name} - <a id="message-header-sign" href="#">{this.props.friends[this.props.target].sign}</a>
        </div>
        <div id="message-content" className="area-content" ref="refContent" onScroll={this.handleScroll}>
        {
          this.props.messages.map(function(msg) {
            if (msg.from == 'system') {
              return <p className={"wrapper-message-" + msg.from}><span className={"message-balloon message-" + msg.from}>{'【' + msg.content + '】'}</span></p>
            }
			return <MessageBalloon msg={msg} />
            // return <p className={"wrapper-message-" + msg.from}><span className={"message-balloon message-" + msg.from}>{msg.content}</span></p>
          })
        }
        </div>
        <div id="message-control-panel" ref="panel">
          <div id="message-box">
            <div id="wrapper-message-box" className="wrapper-input">
              <textarea ref="refInput" type="text" name="id" id="login-id" onKeyPress={this.sendMessageByKeyboard} value={this.state.userInput} onChange={this.handleChange} placeholder="請在這裡輸入訊息！"></textarea>
            </div>
          </div>
          {function() {
            switch (this.props.target) {
              case 0:
                return (
                <div className="pull-left">
                  <a id="button-bonbon" className="message-button" onClick={this.bonbon}>Bonbon!</a>
                  <a id="button-report" className="message-button" onClick={this.disconnect}>離開</a>
                </div>);
              //default:
                //return ();
            }
          }.bind(this)()}
          <div className="pull-right">
            <a id="button-send-image" className="message-button" onclick="return false">傳送圖片</a>
            <a id="button-send-message" className="message-button" onClick={this.sendMessage}>傳送訊息</a>
          </div>
          <div style={{clear: "both"}}>
          </div>
        </div>
      </div>
    );
  }
});

var Chat = React.createClass({
  getInitialState: function() {
    this.props.chatSocket.addHandler('init', function(cmd) {
      var friends = [];
      var initFriend = {
        index: 0,
        name: '陌生人',
        ID: 0,
        online: false,
        stat: 'read',
        img: '/static/img/stranger-m.jpg',
        sign: '猜猜我是誰',
        read: (Date.now() * 10e+5).toString(),
        messages: [{from: 'system', content: '尚未配對成功', time: (Date.now() * 10e+5).toString()}]
      };
      friends.push(initFriend);
      // BUG cmd.Friends may be null
      for (var i = 0; i < cmd.Friends.length; i++) {
        console.log(cmd.Friends[i]);
        var friend = {
          index: i + 1,
          name: cmd.Friends[i].Nick,
          ID: cmd.Friends[i].ID,
          online: cmd.Friends[i].Status == 'on' ? true : false,
          stat: i == 0 ? 'selected' : 'read',
          img: '/static/img/friend_' + parseInt(i + 1) + '.jpg',
          sign: cmd.Friends[i].Sign,
          read: cmd.Friends[i].Read,
          messages: [],
        };
        console.log('read' + friend.read);
        friends.push(friend);
      }
      this.setState({
        friend_number: cmd.Friends.length,
        friends: friends,
        who: 0
      });
      for (var i = 1; i < this.state.friends.length; i++) {
        this.props.chatSocket.send(JSON.stringify({Cmd: "history", With_who: this.state.friends[i].ID, Number: 15, When: (Date.now() * 10e+5).toString(), Order: 0}));
      }
    }.bind(this));

    this.props.chatSocket.addHandler('send', function(cmd) {
      /* send message to sb. */
      var index = -1;
      for (var i = 0; i < this.state.friends.length; i++) {
        if (this.state.friends[i].ID == cmd.Who) {
          index = i;
        }
      }
      if (cmd.OK == true) {
        this.state.friends[index].messages.push({content: cmd.Msg, from: 'me', time: cmd.Time});
      }
      else {
        console.log('send fail(cmd.OK=false)');
      }
      this.setState({
        friends: this.state.friends
      });
    }.bind(this));

    this.props.chatSocket.addHandler('sendFromServer', function(cmd) {
      var index = -1;
      for (var i = 0; i < this.state.friends.length; i++) {
        if (this.state.friends[i].ID == cmd.Who) {
          index = i;
        }
      }
      if (this.state.who != index) {
        NewMessage(this.state.friends[index].name, cmd.Msg);
      }
      this.state.friends[index].messages.push({content: cmd.Msg, from: 'others', time: cmd.Time});
      this.setState({
        friends: this.state.friends
      });
    }.bind(this));

    this.props.chatSocket.addHandler('status', function(cmd) {
      var index = -1;
      for (var i = 0; i < this.state.friends.length; i++) {
        if (this.state.friends[i].ID == cmd.Who) {
          index = i;
        }
      }
      this.state.friends[index].online = (cmd.Status == 'on') ? true : false;
      this.setState({
        friends: this.state.friends
      });
    }.bind(this));

    this.props.chatSocket.addHandler('connect', function(cmd) {
      var friends = this.state.friends;
      friends[0].messages = [{from: 'system', content: '建立配對中...請稍候', time: (Date.now() * 10e+5).toString()}];
      this.setState({
        friends: friends,
        who: 0
      });
    }.bind(this));
    this.props.chatSocket.addHandler('connected', function(cmd) {
      var friends = this.state.friends;
      friends[0].messages = [{from: 'system', content: '已建立新配對，可以開始聊天囉！', time: (Date.now() * 10e+5).toString()}];
      friends[0].online = true;
      friends[0].sign = cmd.Sign;
      this.setState({
        friends: friends,
        who: 0,
      });
    }.bind(this));
    this.props.chatSocket.addHandler('disconnect', function(cmd) {
      this.state.friends[0].messages.push({from: 'system', content: '連線已中斷', time: (Date.now() * 10e+5).toString()});
      this.state.friends[0].online = false;
      this.setState({
        friends: this.state.friends
      });
    }.bind(this));
    this.props.chatSocket.addHandler('disconnected', function(cmd) {
      this.state.friends[0].messages.push({from: 'system', content: '對方已下線，連線中斷', time: (Date.now() * 10e+5).toString()});
      this.state.friends[0].online = false;
      this.setState({
        friends: this.state.friends
      });
    }.bind(this));
    this.props.chatSocket.addHandler('bonbon', function(cmd) {

    }.bind(this));
    this.props.chatSocket.addHandler('new_friend', function(cmd) {
      var index = this.state.friends.length;
      var new_friend = {
        index: index,
        name: cmd.Nick,
        ID: cmd.Who,
        online: true,
        stat: 'selected',
        img: '/static/img/friend_' + index + '.jpg',
        sign: this.state.friends[0].sign,
        read: this.state.friends[0].read,
        messages: this.state.friends[0].messages
      };
      new_friend.messages.push({from: 'system', content: '你們已經Bon在一起，成為了好友！', time: (Date.now() * 10e+5).toString()});
      this.state.friends.push(new_friend);
      this.state.friends[0] = {
        index: 0,
        name: '陌生人',
        ID: 0,
        online: false,
        stat: 'read',
        img: '/static/img/stranger-m.jpg',
        sign: '猜猜我是誰',
        read: cmd.Friends[i].Read,
        messages: [{from: 'system', content: '尚未配對成功', time: (Date.now() * 10e+5).toString()}]
      };
      this.props.chatSocket.addHandler('read', function(cmd) {
        var index = -1;
        for (var i = 0; i < this.state.friends.length; i++) {
          if (this.state.friends[i].ID == cmd.Who) {
            index = i;
          }
        }
        if (index == -1) return 0;
        if (this.state.friends[index].time < cmd.Time) {
          this.state.friends[index].time = cmd.Time;
        }
        this.setState({
          friends: this.state.friends
        });
      }.bind(this));
      this.setState({
        friends: this.state.friends,
        who: index,
      });
    }.bind(this));
    this.props.chatSocket.addHandler('history', function(cmd) {
      var index = -1;
      for (var i = 0; i < this.state.friends.length; i++) {
        if (this.state.friends[i].ID == cmd.With_who) {
          index = i;
        }
      }
      var historyMessage = cmd.Msgs;
      for (var i = 0; i < historyMessage.length; i++) {
        var msg = historyMessage[i];
        if (this.state.friends[index].ID == msg.From) {
          this.state.friends[index].messages.unshift({content: msg.Text, from: 'others', time: msg.Time});
        }
        else {
          this.state.friends[index].messages.unshift({content: msg.Text, from: 'me', time: msg.Time});
        }
      }
      this.setState({
        friends: this.state.friends
      });
    }.bind(this));
    return {
      who: 0,
      friends: [{index: 0,
      name: '陌生人',
      ID: 0,
      online: false,
      stat: 'read',
      img: '/static/img/stranger-m.jpg',
      sign: '猜猜我是誰',
      messages: [{from: 'system', content: '尚未配對成功', time: (Date.now() * 10e+5).toString()}]}],
    };
  },

  selectFriend: function(selectedFriend) {
    this.state.friends[this.state.who].stat = 'read';
    this.state.friends[selectedFriend].stat = 'selected';
    this.setState({
      who: selectedFriend,
      friends: this.state.friends,
    });
    this.refs.refChat.focusInput();
  },

  addMessage: function(who, where, message) {
    if (where == 'buttom') {
      console.log('sending to id: ' + this.state.friends[who].ID);
      try {
        if (this.props.chatSocket.readyState == 1) {
          this.props.chatSocket.send(JSON.stringify({Cmd: "send", Who: this.state.friends[who].ID, Msg: message.content}));
        }
        else {
          console.log('send fail(socket close)');
          return false;
        }
      }
      catch(err) {
        console.log('send fail(socket error)');
        return false;
      }
    }
    return true;
  },

  render: function() {
    return (
    <div id="chat-panel">
      <FriendList friends={this.state.friends} changeState={this.props.changeState} selectedFriend={this.state.who} select={this.selectFriend} chatSocket={this.props.chatSocket}/>
      {function() {
        if (this.props.show == 'chat') {
          return (
            <ChatRoom ref="refChat" messages={this.state.friends[this.state.who].messages} friends={this.state.friends} target={this.state.who} addMessage={this.addMessage} chatSocket={this.props.chatSocket}/>
          );
        }
      }.bind(this)()}
    </div>);
  }
});

//not using but can be used for NewPage
/*
AnotherPage = React.createClass({
  getInitialState: function() {
    return {
    };
  },
  render: function() {

  }
});
*/
var Content = React.createClass({
  render: function() {
    return (
      <Chat chatSocket={this.props.chatSocket} show={this.props.show} changeState={this.props.changeState}/>
    );
  }
});

var ChatPage = React.createClass({
  getInitialState: function() {
    return {
      chatSocket: createSocket(this.props.token),
      show: 'chat',
    };
  },

  componentDidMount: function() {
  },

  componentWillUnmount: function() {
  },

  changeState: function(str) {
    this.setState({
      show: str
    });
  },

  render: function() {
    return (
      <div id="chat-page">
        <SideBar show={this.state.show} changeState={this.changeState} chatSocket={this.state.chatSocket} logout={this.props.logout}/>
        <Content show={this.state.show} changeState={this.changeState} chatSocket={this.state.chatSocket}/>
      </div>
    );
  }
});

module.exports = ChatPage;
