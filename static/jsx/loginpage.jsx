var React = require('../bower/react/react-with-addons.js');

var LoginPage = React.createClass({
  componentDidMount: function() {
    // define Facebook app data and try auto login
    window.fbAsyncInit = function() {
      FB.init({
        appId  : '915780538494020',
        xfbml  : true,
        version: 'v2.4'
      });
      if (localStorage.getItem('login') == 'true') {
        FB.getLoginStatus(function(response) {
          this.statusChangeCallback(response);
        }.bind(this));
      }
    }.bind(this);

    // load Facebook SDK
    (function(d, s, id) {
      var js, fjs = d.getElementsByTagName(s)[0];
      if (d.getElementById(id)) return;
      js = d.createElement(s); js.id = id;
      js.src = "//connect.facebook.net/en_US/sdk.js";
      fjs.parentNode.insertBefore(js, fjs);
    } (document, 'script', 'facebook-jssdk'));
  },

  statusChangeCallback: function(response) {
    if (response.status === 'connected') {
      // update access token and show chat page
      localStorage.setItem('login', 'true');
      this.props.logined(response.authResponse.accessToken);
    }
  },

  checkLoginState: function() {
    FB.getLoginStatus(function(response) {
      this.statusChangeCallback(response);
    }.bind(this));
  },

  handleClick: function() {
    console.log('try login');
    FB.login(this.checkLoginState(), {scope: 'public_profile,email,user_friends'});
  },

  render: function() {
    return (
      <div id="login-page">
        <img src="/static/img/bonbon.png" id="bonbon"></img>
        <a id="login-button" onClick={this.handleClick}>Login</a>
      </div>
    );
  }
});

module.exports = LoginPage;
