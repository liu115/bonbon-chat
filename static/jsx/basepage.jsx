PageAll = React.createClass({
  getInitialState: function() {
    return {
      login: false
    };
  },

  logined: function(token) {
    this.setState({
      login: true,
      token: token
    });
  },

  render: function() {
    if (this.state.login == true) {
      return (
        <ChatPage token={this.state.token}/>
      );
    }
    else {
      return (
        <LoginPage logined={this.logined}/>
      );
    }
  }
});

React.render(
  <PageAll/>,
  document.getElementsByTagName('body')[0]
);
