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
        <div>
          <input type="text" id="sign-input" ref="refInput" value={this.state.value} onKeyPress={this.handleType} onChange={this.handleChange} placeholder="按Enter確認更改簽名"/>
        </div>
      );
    }
    else {
      return (
        <div>
        <a className="profile-status" onClick={this.handleClick}>
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
    return {Sign: "我建超世志，必至無上道"};
  },

  render: function() {
    return (
      //<!-- start of navigation area -->
      <nav id="nav">
        <div id="nav-profile">
          <span className="profile-avatar"><a><img src="img/me_finn.jpg"/></a></span>
          <SignClass sign={this.state.Sign} chatSocket={this.props.chatSocket}/>

        </div>
        <a id="new-connection" onClick={this.props.changeState.bind(null, 'new_connection')}>建立新連線</a>
        <ul id="menu">
          <li><a><span><i className="fa fa-sign-out"></i><span style={{margin: '0px'}}>登出</span></span></a></li>
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
          <p className="friend-info-status">最後的聊天內容</p>
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
    };
  },

  render: function() {
    var friendBoxs = [];
    for (var i = 0; i < this.props.friends.length; i++) {
      friendBoxs.push(<FriendBox index={i} friend={this.props.friends[i]} changeState={this.props.changeState} select={this.props.select}/>);
    }
    return (
      <div id="friend-area">
        <div id="friend-search">
          <div id="wrapper-input-search" className="wrapper-input">
            <input type="text" placeholder="搜尋朋友"/>
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
    return {
      userInput: '',
      scroll: 0,
      roomWidth: window.innerWidth - 521,
      roomHeight: window.innerHeight - 51 - 91 - 15
    };
  },

  handleChange: function(e) {
    this.setState({
      userInput: e.target.value
    });
  },

  sendMessage: function(e) {
    //send it to websocket
    //this.state.messages.splice(0, 0, ['me', 'lalala']);
    if (this.state.userInput != ''){
      this.props.addMessage(this.props.target, 'buttom', {
        from: 'me',
        content: this.state.userInput
      });
      this.setState({
        userInput: ''
      });
      //scrollTop = scrollHeight
    }
    this.setState({
      scroll: React.findDOMNode(this.refs.refContent).scrollHeight
    }, function() {
      React.findDOMNode(this.refs.refContent).scrollTop = (this.state.scroll);
    });
    this.focusInput();
  },

  sendMessageByKeyboard: function(e) {
    var keyInput = e.keyCode == 0 ? e.which : e.keyCode;
    if (keyInput == 13) {
      this.sendMessage();
    }
  },

  focusInput: function() {
    React.findDOMNode(this.refs.refInput).focus();

  },

  handleResize: function(e) {
    this.setState({
      roomWidth: window.innerWidth - 521,
      roomHeight: window.innerHeight - React.findDOMNode(this.refs.header).offsetHeight - React.findDOMNode(this.refs.panel).offsetHeight - 15
    });
  },

  handleScroll: function() {
    this.setState({
      scroll: React.findDOMNode(this.refs.refContent).scrollTop
    });
  },

  componentDidMount: function() {
    window.addEventListener('scroll', this.handleScroll);
    React.findDOMNode(this.refs.refInput).focus();
  },

  componentWillUnmount: function() {
    window.removeEventListener('scroll', this.handleScroll);
  },

  render: function() {
    return (
      <div id="message-area" style={{width: (this.props.roomSize.width - 320 + 'px'), height: (this.props.roomSize.height + 'px')}}>
        <div id="message-header" className="area-header"r ref="header">
          {this.props.friends[this.props.target].name} - <a id="message-header-sign" href="#">{this.props.header}</a>
        </div>
        <div id="message-content" className="area-content" ref="refContent" style={{height: (this.props.roomSize.height - 51 - 91 - 15 + 'px')}}>
        {
          this.props.messages.map(function(msg) {
            return <p><span className={"message-" + msg.from}>{msg.content}</span></p>
          })
        }
        </div>
        <div id="message-panel" ref="panel">
          <div id="message-box">
            <div id="wrapper-message-box" className="wrapper-input">
              <input ref="refInput" type="text" name="id" id="login-id" onKeyPress={this.sendMessageByKeyboard} value={this.state.userInput} onChange={this.handleChange} placeholder="請在這裡輸入訊息！"/>
            </div>
          </div>
          <div className="pull-left">
            <a id="button-bonbon" className="message-button" onclick="return false">Bonbon!</a>
            <a id="button-report" className="message-button" onclick="return false">離開</a>
          </div>
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
        img: 'img/stranger.png',
        sign: '猜猜我是誰',
        messages: [{from: 'system', content: '尚未配對成功'}]
      };
      friends.push(initFriend);
      // BUG cmd.Friends may be null
      for (var i = 0; i < cmd.Friends.length; i++) {
        console.log(cmd.Friends[i])
        var friend = {
          index: i + 1,
          name: cmd.Friends[i].Nick,
          ID: cmd.Friends[i].ID,
          online: cmd.Friends[i].Status == 'on' ? true : false,
          stat: i == 0 ? 'selected' : 'read',
          img: 'img/friend_' + parseInt(i + 1) + '.jpg',
          sign: cmd.Friends[i].Sign,
          messages: [],
        };
        friends.push(friend);
      }
      this.setState({
        friends: friends,
        header: friends[this.state.who].sign,
        who: 1
      });
    }.bind(this));

    this.props.chatSocket.addHandler('send', function(cmd) {
      console.log("something sent!");
      /* send message to sb. */
      var index = -1;
      for (var i = 0; i < this.state.friends.length; i++) {
        if (this.state.friends[i].ID == cmd.Who) {
          index = i;
        }
      }
      console.log('index is ' + index);
      if (cmd.OK == true) {
        this.state.friends[index].messages.push({content: cmd.Msg, from: 'me'});
      }
      else {
        this.state.friends[index].messages.push({content: cmd.Msg + '(send failed)', from: 'me'});
        console.log(this.state.friends[index].messages);
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
      this.state.friends[index].messages.push({content: cmd.Msg, from: 'others'});
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
      friends[0].messages = [{from: 'system', content: '建立配對中...請稍候'}];
      this.setState({
        friends: friends
      });
    }.bind(this));
    this.props.chatSocket.addHandler('connected', function(cmd) {
      var friends = this.state.friends;
      friends[0].messages = [{from: 'system', content: '已建立新配對，可以開始聊天囉！'}];
      friends[0].online = true;
      this.setState({
        friends: friends
      });
    }.bind(this));
    this.props.chatSocket.addHandler('disconnect', function(cmd) {
      this.state.friends[0].messages.push({from: 'system', content: '對方以下線，連線中斷'});
      friends[0].online = false;
      this.setState({
        friends: this.state.friends
      });
    }.bind(this));
    this.props.chatSocket.addHandler('disconnected', function(cmd) {
      this.state.friends[0].messages.push({from: 'system', content: '對方以下線，連線中斷'});
      friends[0].online = false;
      this.setState({
        friends: this.state.friends
      });
    }.bind(this));
    return {
      who: 0,
      friends: [{messages: []}],
      header: '', /* header need fix */
    };
  },

  selectFriend: function(selectedFriend) {
    this.state.friends[this.state.who].stat = 'read';
    this.state.friends[selectedFriend].stat = 'selected';
    this.setState({
      who: selectedFriend,
      friends: this.state.friends,
      header: this.state.friends[selectedFriend].sign
    });
    this.refs.refChat.focusInput();
  },

  addMessage: function(who, where, message) {
    if (where == 'buttom') {
      this.props.chatSocket.send(JSON.stringify({Cmd: "send", Who: this.state.friends[who].ID, Msg: message.content}));
    }

  },

  render: function() {
    if (this.props.show == 'chat') {
      return (
        <div>
          <FriendList friends={this.state.friends} changeState={this.props.changeState} selectedFriend={this.state.who} select={this.selectFriend} chatSocket={this.props.chatSocket}/>
          <ChatRoom ref="refChat" messages={this.state.friends[this.state.who].messages} friends={this.state.friends} target={this.state.who} header={this.state.header} addMessage={this.addMessage} roomSize={this.props.roomSize}/>
        </div>
      );
    }
    else if (this.props.show == 'new_connection') {
      return (
        <div>
          <FriendList friends={this.state.friends} changeState={this.props.changeState} selectedFriend={this.state.who} select={this.selectFriend} chatSocket={this.props.chatSocket}/>
          <NewConnection chatSocket={this.props.chatSocket} changeState={this.props.changeState} roomSize={this.props.roomSize}/>
        </div>
      );
    }
  }
});

var NewConnection = React.createClass({
  getInitialState: function() {
    //name is this.props.name and header take from the name
    return {
    };
  },

  L1Friend: function() {
    this.props.chatSocket.send(JSON.stringify({Cmd: "connect", Type: "L1_FB_friend"}));
    this.handleClick();
  },

  L2Friend: function() {
    this.props.chatSocket.send(JSON.stringify({Cmd: "connect", Type: "L2_FB_friend"}));
    this.handleClick();
  },

  Stranger: function() {
    this.props.chatSocket.send(JSON.stringify({Cmd: "connect", Type: "stranger"}));
    this.handleClick();
  },

  handleClick: function() {
    this.props.changeState('chat');
  },

  render: function() {
    return (
      <div id="connection" style={{width: this.props.roomSize.width - 320 + 'px', height: this.props.roomSize.height + 'px'}}>
        <ul id="connection-list">
          <li className="connection-line"><a className="connection-button" Click={this.L1Friend}>FB的好友</a></li>
          <li className="connection-line"><a className="connection-button" onClick={this.L2Friend}>朋友的朋友</a></li>
          <li className="connection-line"><a className="connection-button" onClick={this.Stranger}>陌生人</a></li>
          <li className="connection-line"><a className="connection-button" onClick={this.handleClick}>取消</a></li>
        </ul>
      </div>
    );
  }
});

var Content = React.createClass({
  render: function() {
    return (
      <Chat chatSocket={this.props.chatSocket} show={this.props.show} changeState={this.props.changeState} roomSize={this.props.roomSize}/>
    );
  }
});

var App = React.createClass({
  getInitialState: function() {
    return {
      chatSocket: createSocket(this.props.token),
      show: 'chat',
      roomWidth: window.innerWidth - 200,
      roomHeight: window.innerHeight
    };
  },

  handleResize: function(e) {
    this.setState({
      roomWidth: window.innerWidth - 200,
      roomHeight: window.innerHeight
    });
  },

  componentDidMount: function() {
    window.addEventListener('resize', this.handleResize);
  },

  componentWillUnmount: function() {
    window.removeEventListener('resize', this.handleResize);
  },

  changeState: function(str) {
    this.setState({
      show: str
    });
  },

  render: function() {
    var size = {width: this.state.roomWidth, height: this.state.roomHeight};
    return (
      <div>
        <SideBar show={this.state.show} changeState={this.changeState} chatSocket={this.state.chatSocket}/>
        <Content show={this.state.show} changeState={this.changeState} chatSocket={this.state.chatSocket} roomSize={size}/>
      </div>
    );
  }
});
