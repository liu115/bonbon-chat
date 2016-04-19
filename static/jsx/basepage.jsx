var React = require('../bower/react/react-with-addons.js');
var ChatPage = require('./chatpage.jsx');
var LoginPage = require('./loginpage.jsx');

if (localStorage.getItem('login') === null) {
  localStorage.setItem('login', 'false');
}

var BasePage = React.createClass({
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
