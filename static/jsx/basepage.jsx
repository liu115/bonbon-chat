if (localStorage.getItem('login') === null) {
  localStorage.setItem('login', 'false');
}
var ReactCSSTransitionGroup = React.addons.CSSTransitionGroup;

BasePage = React.createClass({
  getInitialState: function() {
    return {
      login: false
    };
  },
  logout: function() {
    this.setState({
      login: false
    });
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
        <ChatPage token={this.state.token} logout={this.logout}/>
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
  <BasePage/>,
  document.getElementById('all')
);
