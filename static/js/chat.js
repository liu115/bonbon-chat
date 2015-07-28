//var mySocket = new WebSocket("ws://localhost:8080/chat");
var SideBar = React.createClass({
  getInitialState: function() {
    this.props.chatSocket.addHandler("init", function(cmd) {
      this.setState({Sign: cmd.Setting.Sign})
    }.bind(this))
    return {Sign: "我建超世志，必至無上道"}
  },
  render: function() {
    return (
      //<!-- start of navigation area -->
      <nav id="nav">
        <div id="nav-profile">
          <span className="profile-avatar"><a><img src="img/me_finn.jpg"/></a></span>
          <a className="profile-status">{this.state.Sign}</a>
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
    this.props.select(this.props.friend.index);
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
    return {
    };
  },
  render: function() {
    return (
      <div id="friend-area">
        <div id="friend-search">
          <div id="wrapper-input-search" className="wrapper-input">
            <input type="text" placeholder="搜尋朋友"/>
          </div>
        </div>
        {
          this.props.friends.map(function(friend){
            return (
              <FriendBox friend={friend} select={this.props.select}/>
            );
          }.bind(this))
        }
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
    React.findDOMNode(this.refs.refInput).focus();
  },
  sendMessageByKeyboard: function(e) {
    var keyInput = e.keyCode == 0 ? e.which : e.keyCode;
    if (keyInput == 13) {
      this.sendMessage();
    }
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


var Content = React.createClass({
  getInitialState: function() {
    return {
      who: 1,
      header: 'Its where my demons hide.', /* header need fix */
      messages: [
        [
          {
            from: 'system',
            content: '已建立連線，開始聊天吧！'
          }
        ],
        [
          {
            from: 'system',
            content: '已建立連線，開始聊天吧！'
          },
          {
            from: 'me',
            content: 'hihi'
          },
          {
            from: 'others',
            content: 'abcde'
          }
        ],
        [
          {
            from: 'system',
            content: '已建立連線，開始聊天吧！'
          }
        ],
        [],
        [],
        []
      ],
      friends: [
        {
          index: 0,
          name: '陌生人',
          stat: 'read',
          online: true,
          img: 'img/friend_0.jpg'
        },
        {
          index: 1,
          name: 'Apple',
          stat: 'selected',
          online: true,
          img: 'img/friend_1.jpg'
        },
        {
          index: 2,
          name: 'Banana',
          stat: 'read',
          online: true,
          img: 'img/friend_2.jpg'
        },
        {
          index: 3,
          name: 'Cake',
          stat: 'unread',
          online: true,
          img: 'img/friend_3.jpg'
        },
        {
          index: 4,
          name: 'Donut',
          stat: 'read',
          online: false,
          img: 'img/friend_4.jpg'
        },
        {
          index: 5,
          name: 'Egg',
          stat: 'unread',
          online: false,
          img: 'img/friend_5.jpg'
        }
      ]
    };
  },
  selectFriend: function(selectedFriend) {
    this.state.friends[this.state.who].stat = 'read';
    this.state.friends[selectedFriend].stat = 'selected';
    this.setState({
      who: selectedFriend,
      friends: this.state.friends
      //set header from data base
    });
  },
  addMessage: function(who, where, message) {
    if (where == 'buttom') {
      this.state.messages[who].push(message);
    }
    this.setState({
      messages: this.state.messages
    });
  },
  render: function() {
    return (
      <div>
        <FriendList friends={this.state.friends} selectedFriend={this.state.who} select={this.selectFriend}/>
        <ChatRoom messages={this.state.messages[this.state.who]} target={this.state.who} friends={this.state.friends} header={this.state.header} addMessage={this.addMessage}/>
      </div>
    );
  }
});
var App = React.createClass({
  getInitialState: function() {
    return {chatSocket: createSocket()}
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
