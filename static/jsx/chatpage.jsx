var React = require('../bower/react/react-with-addons.js');
var MessageBalloon  = require('./message.jsx');
var ReactCSSTransitionGroup = React.addons.CSSTransitionGroup;

var SignClass = React.createClass({
  getInitialState: function() {
    return {
      setting: false,
      value: this.props.sign
    };
  },

  componentDidUpdate: function (prevProps, prevState) {
    if (!prevState.setting && this.state.setting) {
      var input = React.findDOMNode(this.refs.refInput)
      input.focus();
      input.selectionStart = input.selectionEnd = this.state.value.length;
    }
  },

  componentWillReceiveProps: function(nextProps) {
    if (nextProps.sign != this.props.sign) {
      this.setState({value: nextProps.sign})
    }
  },

  handleClick: function() {
    this.setState({setting: true});
  },

  handleType: function(e) {
    var keyInput = e.keyCode == 0 ? e.which : e.keyCode;
    if (keyInput == 13) {
      console.log('trying to set ' + this.state.value + ' as Sign.');
      this.props.chatSocket.send(JSON.stringify({Cmd: "setting", Setting: {Sign: this.state.value}}));
      this.setState({
        setting: false,
        value: e.target.value
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
	          <input type="text" id="sign-input" ref="refInput" value={this.state.value} onKeyPress={this.handleType} onChange={this.handleChange} placeholder="按Enter確認更改"/>
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
      this.setState({
        Sign: cmd.Setting.Sign,
        Avatar: cmd.Setting.Avatar,
      });
    }.bind(this));
    this.props.chatSocket.addHandler("setting", function(cmd) {
      if (cmd.OK == true) {
        if (cmd.Setting.Sign) {
          console.log('setting success, new sign is ' + cmd.Setting.Sign);
          this.setState({Sign: cmd.Setting.Sign});
        }
        if (cmd.Setting.Avatar) {
          this.setState({Avatar: cmd.Setting.Avatar});
          console.log('setting success, new avatar is ' + cmd.Setting.Avatar);
        }
      }
      else {
        console.log('setting failed!');
      }
    }.bind(this));
    return {Sign: "", Avatar: "", buttonList: null, selectAvatar: false};
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
  handleStartSelectAvatar: function(e) {
    this.setState({selectAvatar: true})
    e.preventDefault();
  },
  handleEndAvatar: function(e) {
    this.setState({selectAvatar: false})
    e.preventDefault();
  },
  handleSelectAvatar: function(avatar) {
    return function () {
      this.props.chatSocket.send(JSON.stringify({Cmd: "setting", Setting: {Avatar: avatar}}));
    }.bind(this)
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
      // <!-- start of navigation area -->
      <nav id="sidebar-panel">
        <div id="sidebar-profile">
          <div id="profile-avatar">
            <img id="my-avatar" src={"/static/img/avatar/" + this.state.Avatar + ".jpg"} />
            <a href="#" onClick={this.handleStartSelectAvatar}>
              <i className="fa fa-user"></i>
              <span className="change-avatar-text">點我改大頭</span>
            </a>
            {function() {
              if (this.state.selectAvatar) {
                var avatars = ['換個大頭貼吧', 'ㄇㄐ', '阿砲', '毛毛', '花惹發', '桃子'];
                return (
                  <div id="avatar-list-wrap">
                    <div id="end-avatar-list" onClick={this.handleEndAvatar}>
                      <i className="fa fa-times"></i>
                    </div>
                    <div id="avatar-list">
                      {
                        avatars.map(
                          function(a) {
                            return (
                              <div data-balloon={a} data-balloon-pos="down"
                                className="avatar-to-select" onClick={this.handleSelectAvatar(a)}>
                                <img src={"/static/img/avatar/" + a + ".jpg"}/>
                              </div>
                              )
                          }.bind(this)
                        )
                      }
                    </div>
                  </div>
                )
              }
              return null;
            }.bind(this)()}
          </div>
          <SignClass sign={this.state.Sign} chatSocket={this.props.chatSocket}/>
        </div>
        <a id="new-connection" onClick={this.handleClick}>建立新連線</a>
        <ReactCSSTransitionGroup transitionName="list-animate">
          {this.state.buttonList}
        </ReactCSSTransitionGroup>
        <ul id="menu">
          <li><a onClick={this.logout}><span><i className="fa fa-sign-out"></i><span style={{margin: '0px'}}>登出</span></span></a></li>
          {/*TOTO: 應建立一份前端設設定檔，將開發者信箱置於此設定檔中*/}
          <li><a href="mailto:yc1043@gmail.com"><span><i className="fa fa-envelope"></i><span style={{margin: '0px'}}>聯絡開發者</span></span></a></li>
        </ul>
      </nav>
      //<!-- end of navigation area -->
    );
  }
});

var FriendBox = React.createClass({
  getInitialState: function () {
    return {
      name: this.props.friend.name,
      changing: false,
    }
  },
  componentWillReceiveProps: function(nextProps) {
    if (nextProps.friend.name != this.props.friend.name) {
      this.setState({
        changing: false,
        name: nextProps.friend.name
      });
    }
  },
  handleClick: function() {
    this.props.select(this.props.index);
    this.props.changeState('chat');
  },
  handleClickPencil: function (e) {
    this.setState({changing: true});
    e.stopPropagation();
  },
  handleType: function (e) {
    var keyInput = e.keyCode == 0 ? e.which : e.keyCode;
    if (keyInput == 13) {
      if (e.target.value != this.props.friend.name) {
        console.log("try set_nick of id " + this.props.friend.ID + " to " + e.target.value);
        this.props.chatSocket.send(JSON.stringify({Cmd: "set_nick", Who: this.props.friend.ID, Nick: e.target.value}));
      } else {
        this.setState({changing: false})
      }
    }
  },
  handleChange: function (e) {
    this.setState({name: e.target.value})
  },
  handleCheck: function(e) {
    if (this.state.name != this.props.friend.name) {
      console.log("try set_nick of id " + this.props.friend.ID + " to " + this.state.name);
      this.props.chatSocket.send(JSON.stringify({Cmd: "set_nick", Who: this.props.friend.ID, Nick: this.state.name}));
    } else {
      this.setState({changing: false})
    }
    e.stopPropagation();
  },
  handleCancel: function() {
    this.setState({
      changing: false,
      name: this.props.friend.name
    })
    e.stopPropagation();
  },
  componentDidUpdate: function (prevProps, prevState) {
    if (!prevState.changing && this.state.changing) {
      var input = this.refs.InputName.getDOMNode();
      input.focus();
      input.selectionStart = input.selectionEnd = this.state.name.length;
    }
  },
  render: function() {
    return (
      <div className={"friend-unit " + "friend-" + this.props.friend.stat + (this.props.friend.online ? '' : " off-line")}
           onClick={this.handleClick}>
        <div className={(this.props.index == 0) ? "stranger-avatar": "friend-avatar"}>
          <img src={this.props.friend.img}/>
        </div>
        <div className="friend-info">
            {function () {
              if (this.state.changing) {
                return <input type="text" value={this.state.name}
                  ref="InputName" onChange={this.handleChange}
                  onClick={function(e) {e.stopPropagation()}}
                  onKeyPress={this.handleType}/>
              } else {
                return <p className="friend-info-name"> {this.props.friend.name} </p>
              }
            }.bind(this)()}
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
        <div className="friend-setting">
          {function () {
            if (this.state.changing) {
              return [<i className="fa fa-check" onClick={this.handleCheck}></i>,
                <i className="fa fa-times" onClick={this.handleCancel}></i>];
            } else {
              return <i className="fa fa-pencil" onClick={this.handleClickPencil}></i>;
            }
          }.bind(this)()}
        </div>
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
      friendBoxs.push(<FriendBox chatSocket={this.props.chatSocket} index={i} key={i} friend={this.props.friends[i]} changeState={this.props.changeState} select={this.props.select}/>);
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

var InputArea = React.createClass({
	getInitialState: function () {
		/* lock when message are sending */
		return {
		};
  },
  sendMessageByKeyboard: function(e) {
    var keyInput = e.keyCode == 0 ? e.which : e.keyCode;
    if (keyInput == 13 && !e.shiftKey) {
      this.props.sendMessage();
      e.preventDefault();
    }
  },
  handleChange: function(e) {
    this.props.changeInput(e.target.value);
  },
  focusInput: function() {
    React.findDOMNode(this.refs.refInput).focus();
  },
  render: function () {
    return <textarea ref="refInput" type="text" name="id" id="login-id" onKeyPress={this.sendMessageByKeyboard} value={this.props.userInput} onChange={this.handleChange} placeholder="請在這裡輸入訊息！"></textarea>
  }
})

var ChatRoom = React.createClass({
  getInitialState: function() {
    //name is this.props.name and header take from the name
    this.props.chatSocket.addHandler('history', function(cmd) {
      var node = React.findDOMNode(this.refs.refContent);
      node.scrollTop = this.scrollTop + (node.scrollHeight - this.scrollHeight);
    }.bind(this));
		this.props.chatSocket.addHandler('bonbon', function(cmd) {
			if (cmd.OK == true) this.setState({bonboning: true});
    }.bind(this));
		this.props.chatSocket.addHandler('disconnect', function(cmd) {
			if (cmd.OK == true) this.setState({bonboning: false});
    }.bind(this));
		this.props.chatSocket.addHandler('new_friend', function(cmd) {
			this.setState({bonboning: false});
		}.bind(this));
    this.props.chatSocket.addHandler('send', function(cmd) {
			if (cmd.OK == true) {
				this.setState({
					userInput: ''
				});
        this.freeSendLock();
			}
		}.bind(this));
	return {
		bonboning: false,
    userInput: '',
    sendLock: false
		};
  },
  componentWillUpdate: function() {
    var node = React.findDOMNode(this.refs.refContent);
    this.scrollHeight = node.scrollHeight;
    this.scrollTop = node.scrollTop;
    this.shouldScrollBottom = node.scrollTop + node.offsetHeight === node.scrollHeight;
  },
  componentDidUpdate: function() {
		/* If the scroll is at bottom then send 'read'. */
    if (this.shouldScrollBottom) {
      var node = React.findDOMNode(this.refs.refContent);
      node.scrollTop = node.scrollHeight;
      var messageLength = this.props.friends[this.props.target].messages.length;
      if (messageLength <= 0) return 0;
      var lastMessage = false;
      lastMessage = this.props.friends[this.props.target].messages[messageLength - 1];
      if (lastMessage === false) return 0;
      if (this.props.friends[this.props.target].read < lastMessage.time) {
        this.props.chatSocket.send(JSON.stringify({Cmd: "read", With_who: this.props.friends[this.props.target].ID}));
      }
			this.shouldScrollBottom = false;
    }
  },
  componentDidMount: function() {
    this.historyLock = 'unlock'
    React.findDOMNode(this.refs.refInput).focus();
  },
  componentWillUnmount: function() {
  },
  requireSendLock: function() {
    if (this.state.sendLock) return false;
    this.setState({sendLock: true});
    return true;
  },
  freeSendLock: function() {
    this.setState({sendLock: false});
  },
  sendMessage: function(e) {
    //send it to websocket
    //this.state.messages.splice(0, 0, ['me', 'lalala']);
      this.props.addMessage(this.props.target, 'buttom', {
        from: 'me',
        content: this.state.userInput.trim()
      });
    this.focusInput();
  },
  changeInput: function(value) {
    this.setState({
      userInput: value
    });
  },
	autoScrollBottom: function() {
		if (this.shouldScrollBottom) {
      var node = React.findDOMNode(this.refs.refContent);
      node.scrollTop = node.scrollHeight;
		}
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
		/*
    if (node.scrollTop + node.offsetHeight === node.scrollHeight) {
      var messageLength = this.props.friends[this.props.target].messages.length;
      if (messageLength <= 0) return 0;
      var lastMessage = false;
      //for (var i = messageLength - 1; i >= 0; i--) {
        //if (this.props.friends[this.props.target].messages[i] == 'other')
          lastMessage = this.props.friends[this.props.target].messages[messageLength - 1];
					//break;
      //}
      if (lastMessage === false) return 0;
			console.log(this.props.friends[this.props.target].read);
			console.log(lastMessage.time);
      if (this.props.friends[this.props.target].read < lastMessage.time) {
        console.log('read');
        this.props.chatSocket.send(JSON.stringify({Cmd: "read", With_who: this.props.friends[this.props.target].ID}));
      }
    }
		*/
  },
  focusInput: function() {
    this.refs.refInput.focusInput();
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
						return <MessageBalloon key={msg.time} msg={msg} autoScroll={this.autoScrollBottom}/>
          }.bind(this))
        }
        </div>
        <div id="message-control-panel" ref="panel">
          <div id="message-box">
            <div id="wrapper-message-box" className="wrapper-input">
              <InputArea ref="refInput" userInput={this.state.userInput} changeInput={this.changeInput} sendMessage={this.sendMessage}/>
            </div>
          </div>
          <div className="pull-left">
          {function() {
            switch (this.props.target) {
              case 0:
								switch (this.state.bonboning) {
									case true:
										return (
                      <span>
			                  <a id="button-bonbon" className="message-button bonbon-button-clicked" onClick={this.bonbon}>Bonbon!</a>
			                  <a id="button-report" className="message-button" onClick={this.disconnect}>離開</a>
                      </span>
                    );
									case false:
										return (
                      <span>
			                  <a id="button-bonbon" className="message-button" onClick={this.bonbon}>Bonbon!</a>
			                  <a id="button-report" className="message-button" onClick={this.disconnect}>離開</a>
                      </span>
			              );
								}
            }
          }.bind(this)()}
          {function() {
            if (this.state.sendLock) {
              return (<span className="deliver-message">傳送中...</span>);
            }
          }.bind(this)()}
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
        stat: 'selected',
        img: '/static/img/avatar/stranger-m.jpg',
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
          stat: 'read',
          img: '/static/img/avatar/' + cmd.Friends[i].Avatar + '.jpg',
          sign: cmd.Friends[i].Sign,
          read: cmd.Friends[i].LastRead,
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
			this.sortFriends();
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
      NewMessage(this.state.friends[index].name, cmd.Msg);
			this.state.friends[index].messages.push({content: cmd.Msg, from: 'others', time: cmd.Time});
			// Sorted by
			this.sortFriends();
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
		this.props.chatSocket.addHandler('set_nick', function(cmd) {
			var index = -1;
      for (var i = 0; i < this.state.friends.length; i++) {
        if (this.state.friends[i].ID == cmd.Who) {
          index = i;
        }
      }
			if (index == -1) return 0; // Who not found
      console.log(this.state.friends[index]);
			this.state.friends[index] = React.addons.update(this.state.friends[index], {name: {$set: cmd.Nick}})
			this.setState({
        friends: this.state.friends
      });
		}.bind(this));
		this.props.chatSocket.addHandler('change_sign', function(cmd) {
			var index = -1;
      for (var i = 0; i < this.state.friends.length; i++) {
        if (this.state.friends[i].ID == cmd.Who) {
          index = i;
        }
      }
			if (index == -1) return 0; // Who not found
			this.state.friends[index].sign = cmd.Sign;
			this.setState({
        friends: this.state.friends
      });
		}.bind(this));
		this.props.chatSocket.addHandler('change_avatar', function(cmd) {
			var index = -1;
      for (var i = 0; i < this.state.friends.length; i++) {
        if (this.state.friends[i].ID == cmd.Who) {
          index = i;
        }
      }
			if (index == -1) return 0; // Who not found
      this.state.friends[index].img = '/static/img/avatar/' + cmd.Avatar + '.jpg';
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
      friends[0].img = '/static/img/avatar/' + cmd.Avatar + '.jpg';
      friends[0].online = true;
      friends[0].sign = cmd.Sign;
      this.setState({
        friends: friends,
        who: 0,
      });
    }.bind(this));
    this.props.chatSocket.addHandler('disconnect', function(cmd) {
			if (cmd.OK == true) {
      	this.state.friends[0].messages.push({from: 'system', content: '連線已中斷', time: (Date.now() * 10e+5).toString()});
      	this.state.friends[0].online = false;
			}
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
      NewFriend();
      var index = this.state.friends.length;
      var new_friend = {
        index: index,
        name: cmd.Nick,
        ID: cmd.Who,
        online: true,
        stat: 'selected',
        img: this.state.friends[0].img,
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
        img: '/static/img/avatar/stranger-m.jpg',
        sign: '猜猜我是誰',
        read: (Date.now() * 10e+5).toString(),
        messages: [{from: 'system', content: '尚未配對成功', time: (Date.now() * 10e+5).toString()}]
      };
			this.setState({
	      friends: this.state.friends,
	      who: index,
	    });
		}.bind(this));
    this.props.chatSocket.addHandler('read', function(cmd) {
      var index = -1;
      for (var i = 0; i < this.state.friends.length; i++) {
        if (this.state.friends[i].ID == cmd.Who) {
          index = i;
        }
      }
      if (index == -1) return 0;
      if (this.state.friends[index].read < cmd.Time) {
        this.state.friends[index].read = cmd.Time;
      }
      this.setState({
        friends: this.state.friends
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
			this.sortFriends();
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
      img: '/static/img/avatar/stranger-m.jpg',
      sign: '猜猜我是誰',
			read: (Date.now() * 10e+5).toString(),
      messages: [{from: 'system', content: '尚未配對成功', time: (Date.now() * 10e+5).toString()}]}],
    };
  },
	sortFriends: function() {
		var	who = this.state.friends[this.state.who].ID;
		var sorted_friends = this.state.friends.slice(1, this.state.friends.length).sort(
			function(x , y) {
        var x_time = (x.messages.length == 0) ? 0 : x.messages[x.messages.length - 1].time;
        var y_time = (y.messages.length == 0) ? 0 : y.messages[y.messages.length - 1].time;
        return y_time - x_time;
			}
		);
		sorted_friends.unshift(this.state.friends[0]);
		this.state.friends = sorted_friends;
		for (var i = 0; i < this.state.friends.length; i++) {
			if (this.state.friends[i].ID == who) {
				this.state.who = i;
			}
		}
		console.log("sorting");
		this.setState({
			friends: this.state.friends,
			who: this.state.who
		});
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
