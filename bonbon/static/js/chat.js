var mySocket = new WebSocket("ws://localhost:8080/chat");
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
      <div className="friend-unit friend-read">
        <div className="friend-avatar">
          <img src="img/friend_0.jpg"/>
        </div>
        <div className="friend-info">
          <p className="friend-info-name">陌生人</p>
          <p className="friend-info-status">最後的聊天內容</p>
        </div>
        <div style={{clear: "both"}}></div>
      </div>
    );
  }
});
var FriendList = React.createClass({
  render: function() {
    return (
      <div id="friend-area">
    		<div id="friend-search">
    			<div id="wrapper-input-search" className="wrapper-input">
    				<input type="text" placeholder="搜尋朋友"/>
    			</div>
    		</div>
    		// stranger
    		<div className="friend-unit friend-read">
    			<div className="friend-avatar">
    				<img src="img/friend_0.jpg"/>
    			</div>
    			<div className="friend-info">
    				<p className="friend-info-name">陌生人</p>
    				<p className="friend-info-status">最後的聊天內容</p>
    			</div>
    			<div style={{clear: "both"}}></div>
    		</div>
    		//message: selected
    		<div className="friend-unit friend-selected">
    			<div className="friend-avatar">
    				<img src="img/friend_1.jpg"/>
    			</div>
    			<div className="friend-info">
    				<p className="friend-info-name">Apple</p>
    				<p className="friend-info-status">目前選取的朋友</p>
    			</div>
    			<div style={{clear: "both"}}></div>
    		</div>
    		// message: on-line, have read
    		<div className="friend-unit friend-read">
    			<div className="friend-avatar">
    				<img src="img/friend_2.jpg"/>
    			</div>
    			<div className="friend-info">
    				<p className="friend-info-name">Banana</p>
    				<p className="friend-info-status">上線的朋友、所有訊息已讀</p>
    			</div>
    			<div style={{clear: "both"}}></div>
    		</div>
    		// message: on-line, something unread
    		<div className="friend-unit friend-unread">
    			<div className="friend-avatar">
    				<img src="img/friend_3.jpg"/>
    			</div>
    			<div className="friend-info">
    				<p className="friend-info-name">Cake</p>
    				<p className="friend-info-status">上線的朋友、有訊息未讀</p>
    			</div>
    			<div style={{clear: "both"}}></div>
    		</div>
    		// message: off-line, have unread
    		<div className="friend-unit friend-read off-line">
    			<div className="friend-avatar">
    				<img src="img/friend_4.jpg"/>
    			</div>
    			<div className="friend-info">
    				<p className="friend-info-name">Donut</p>
    				<p className="friend-info-status">離線的朋友、所有訊息已讀</p>
    			</div>
    			<div style={{clear: "both"}}></div>
    		</div>
    		<div className="friend-unit friend-unread off-line">
    			<div className="friend-avatar">
    				<img src="img/friend_5.jpg"/>
    			</div>
    			<div className="friend-info">
    				<p className="friend-info-name">Egg</p>
    				<p className="friend-info-status">離線的朋友、有訊息未讀</p>
    			</div>
    			<div style={{clear: "both"}}></div>
    		</div>
    	</div>
    );
  }
});

var ChatRoom = React.createClass({
  getInitialState: function() {
    //name is this.props.name and header take from the name
    return {
      messages: [
        ['system', '已建立連線，開始聊天吧！'], ['me', 'hihi'], ['others', 'abcde']
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
    this.state.messages.push(['me', 'lalala']);
    this.setState({
      messages: this.state.messages,
      userInput: ''
    });
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
            return <p><span className={"message-" + msg[0]}>{msg[1]}</span></p>
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
