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
          <input type="text" className="sign-input" ref="refInput" value={this.state.value} onKeyPress={this.handleType} onChange={this.handleChange} placeholder="按Enter確認更改簽名"/>
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
        <a id="new-connection">建立新連線</a>
        <ul id="menu">
          <li><a><span><i className="fa fa-comment"></i><span style={{margin: '0px'}}>朋友列表</span></span></a></li>
          <li><a><span><i className="fa fa-cog"></i><span style={{margin: '0px'}}>標籤設定</span></span></a></li>
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
  },
  render: function() {
    return (
      <div className={"friend-unit " + "friend-" + this.props.friend.stat + (this.props.friend.online ? '' : " off-line")} onClick={this.handleClick}>
        <div className="friend-avatar">
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
      friendBoxs.push(<FriendBox index={i} friend={this.props.friends[i]} select={this.props.select}/>);
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
    window.addEventListener('resize', this.handleResize);
    window.addEventListener('scroll', this.handleScroll);
    React.findDOMNode(this.refs.refInput).focus();
  },
  componentWillUnmount: function() {
    window.removeEventListener('resize', this.handleResize);
    window.removeEventListener('scroll', this.handleScroll);
  },
  render: function() {
    return (
      <div id="message-area" style={{width: (this.state.roomWidth + 'px')}}>
        <div id="message-header" ref="header">
          {this.props.friends[this.props.target].name} - <a id="message-header-sign" href="#">{this.props.header}</a>
        </div>
        <div id="message-content" ref="refContent" style={{height: (this.state.roomHeight + 'px')}}>
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
    			<div style={{clear: "both"}}></div>
    		</div>
    	</div>
    );
  }
});


var Chat = React.createClass({
  getInitialState: function() {
    this.props.chatSocket.addHandler('init', function(cmd) {
      var friends = [];
      for (var i = 0; i < cmd.Friends.length; i++) {
        console.log(cmd.Friends[i])
        var friend = {
          index: i,
          name: cmd.Friends[i].Nick,
          ID: cmd.Friends[i].ID,
          online: true,
          stat: i == 0 ? 'selected' : 'read',
          img: 'img/friend_' + parseInt(i + 1) + '.jpg',
          sign: cmd.Friends[i].Sign,
          messages: [],
        };
        friends.push(friend);
      }
      this.setState({
        friends: friends,
        header: friends[this.state.who].sign
      });
    }.bind(this));

    this.props.chatSocket.addHandler('send', function(cmd) {
      console.log("something sent!");
      /* send message to sb. */
      if (cmd.OK == true) {
        this.state.friends[cmd.Who].messages.push(cmd.Msg);
      }
      else {
        this.state.friends[cmd.Who].messages.push(cmd.Msg + '(send failed)');
      }
      this.setState({
        friends: this.state.friends
      });
    }.bind(this));

    this.props.chatSocket.addHandler('sendFromServer', function(cmd) {
      this.state.friends[cmd.Who].messages.push(cmd.Msg);
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
      this.props.chatSocket.send(JSON.stringify({Cmd: "send", Who: who, Msg: message.content}));
    }

  },
  render: function() {
    return (
      <div>
        <FriendList friends={this.state.friends} selectedFriend={this.state.who} select={this.selectFriend} chatSocket={this.props.chatSocket}/>
        <ChatRoom ref="refChat" messages={this.state.friends[this.state.who].messages} friends={this.state.friends} target={this.state.who} header={this.state.header} addMessage={this.addMessage}/>
      </div>
    );
  }
});
var NewConnection = React.createClass({
  render: function() {
    return (
      <div id="connection">
        <ul>
          <li><a>FB的好友</a></li>
          <li><a>朋友的朋友</a></li>
          <li><p>陌生人</p></li>
        </ul>
      </div>
    );
  }
});
var Content = React.createClass({
  getInitialState: function() {
    return {
      show: 'chat'
    };
  },
  changeState: function(str) {
    this.setState({
      show: str
    });
  },
  render: function() {
    if (this.state.show == 'chat') {
      return (
        <Chat chatSocket={this.props.chatSocket} changeState={this.changeState}/>
      );
    }
    if (this.state.show == 'new_connection') {
      return (
        <NewConnection chatSocket={this.props.chatSocket} changeState={this.changeState}/>
      );
    }
  }
});

var App = React.createClass({
  getInitialState: function() {
    return {chatSocket: createSocket(this.props.token)}
  },
  render: function() {
    return (
      <div>
        <SideBar chatSocket={this.state.chatSocket}/>
        <Content chatSocket={this.state.chatSocket}/>
      </div>
    )
  }
})
