//var mySocket = new WebSocket("ws://localhost:8080/chat");
var SideBar = React.createClass({
  render: function() {
    return (
      //<!-- start of navigation area -->
      <nav id="nav">
        <div id="nav-profile">
          <span className="profile-avatar"><a><img src="img/me_finn.jpg"/></a></span>
          <a className="profile-status">這是我的簽名檔（狀態）</a>
        </div>
        <a id="new-connection">建立新連線</a>
        <ul id="menu">
          <li><a><span><i className="fa fa-comment"></i>朋友列表</span></a></li>
          <li><a><span><i className="fa fa-cog"></i>標籤設定</span></a></li>
          <li><a><span><i className="fa fa-sign-out"></i>登出</span></a></li>
        </ul>
      </nav>
      //<!-- end of navigation area -->
    );
  }
});
var FriendBox = React.createClass({
  render: function() {
    return (
      <div>
        <div className="friend-avatar">
          <img src={this.props.img}/>
        </div>
        <div className="friend-info">
          <p className="friend-info-name">{this.props.name}</p>
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
      selected: 1,
      friends: [
        [0, '陌生人', 'read', '', 'img/friend_0.jpg'],
        [1, 'Apple', 'selected', '', 'img/friend_1.jpg'],
        [2, 'Banana', 'read', '', 'img/friend_2.jpg'],
        [3, 'Cake', 'unread', '', 'img/friend_3.jpg'],
        [4, 'Donut', 'read', 'off-line', 'img/friend_4.jpg'],
        [5, 'Egg', 'unread', 'off-line', 'img/friend_5.jpg']
      ]
    };
  },
  friendClickHandler: function(e) {
    alert("clicked");
    this.state.friends[this.state.selected][2] = 'read';
    this.state.selected = e.target.ref.substring(6);
    this.state.friends[this.state.selected][2] = 'selected';
    this.setState({
      selected: this.state.selected,
      friends: this.state.friends
    });
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
          this.state.friends.map(function(friend){
            return (
              <div className={"friend-unit " + "friend-" + friend[2] + " " + friend[3]} onClick={this.friendClickHandler}>
                <FriendBox ref={'friend'+friend[0]} name={friend[1]} img={friend[4]}/>
              </div>
            )
          })
        }
      </div>
    );
  }
});

var ChatRoom = React.createClass({
  getInitialState: function() {
    //name is this.props.name and header take from the name
    return {
      messages: [
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
    this.state.messages.push({
      from: 'me',
      content: this.state.userInput
    });
    this.setState({
      messages: this.state.messages,
      userInput: ''
    });
    }
  },
  handleResize: function(e) {
    this.setState({
      roomWidth: window.innerWidth - 521,
      roomHeight: window.innerHeight - React.findDOMNode(this.refs.header).offsetHeight - React.findDOMNode(this.refs.panel).offsetHeight - 15
    });
  },
  componentDidMount: function() {
    window.addEventListener('resize', this.handleResize);
  },
  componentWillUnmount: function() {
    window.removeEventListener('resize', this.handleResize);
  },
  render: function() {
    return (
      <div id="message-area" style={{width: (this.state.roomWidth + 'px')}}>
        <div id="message-header" ref="header">
          {this.props.name} - <a id="message-header-sign" href="#">{this.props.header}</a>
        </div>
        <div id="message-content" style={{height: (this.state.roomHeight + 'px')}}>
        {
          this.state.messages.map(function(msg) {
            return <p><span className={"message-" + msg.from}>{msg.content}</span></p>
          })
        }
        </div>
        <div id="message-panel" ref="panel">
          <div id="message-box">
            <div id="wrapper-message-box" className="wrapper-input">
              <input type="text" name="id" id="login-id" value={this.state.userInput} onChange={this.handleChange} placeholder="請在這裡輸入訊息！"/>
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
      who: 'Apple',
      header: 'Its where my demons hide.'
    };
  },
  render: function() {
    return (
      <div>
        <FriendList />
        <ChatRoom name={this.state.who} header={this.state.header}/>
      </div>
    );
  }
});

var PageAll = React.createClass({
  render: function() {
    return (
      <div>
        <SideBar/>
        <Content/>
      </div>
    );
  }
});

React.render(
  <div>
    <PageAll/>
  </div>,
  document.getElementById('all')
);
