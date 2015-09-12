LoginPage = React.createClass({
  componentDidMount: function() {
    // define Facebook app data and try auto login
    window.fbAsyncInit = function() {
      FB.init({
        appId  : '915780538494020',
        xfbml  : true,
        version: 'v2.4'
      });
      FB.getLoginStatus(function(response) {
        this.statusChangeCallback(response);
      }.bind(this));
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
      this.props.logined(response.authResponse.accessToken);
    }
  },

  checkLoginState: function() {
    FB.getLoginStatus(function(response) {
      this.statusChangeCallback(response);
    }.bind(this));
  },

  handleClick: function() {
    FB.login(this.checkLoginState(), {scope: 'public_profile,email,user_friends'});
  },

  render: function() {
    return (
      <div id="login-page">
        <div id="login-img">
          <a id="login-button" href="#" onClick={this.handleClick}></a>
        </div>
      </div>
    );
  }
});

window.LoginPage = LoginPage
