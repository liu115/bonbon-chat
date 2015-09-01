var LoginPage = React.createClass({
  componentDidMount: function() {
    window.fbAsyncInit = function() {
      FB.init({
          appId      : '915780538494020',
          xfbml      : true,
          version    : 'v2.4'
      });
      FB.getLoginStatus(function(response) {
        this.statusChangeCallback(response);
      }.bind(this));
    }.bind(this);

    (function(d, s, id) {
      var js, fjs = d.getElementsByTagName(s)[0];
      if (d.getElementById(id)) return;
      js = d.createElement(s); js.id = id;
      js.src = "//connect.facebook.net/en_US/sdk.js";
      fjs.parentNode.insertBefore(js, fjs);
    } (document, 'script', 'facebook-jssdk'));
  },
  testAPI: function() {
    console.log('Welcome!  Fetching your information.... ');
    FB.api('/me', function(response) {
      console.log('Successful login for: ' + response.name);
      //document.getElementById('status').innerHTML =
      //'Thanks for logging in, ' + response.name + '!';
    }.bind(this));
  },
  statusChangeCallback: function(response) {
    console.log('statusChangeCallback');
    console.log(response);
  // The response object is returned with a status field that lets the
  // app know the current login status of the person.
  // Full docs on the response object can be found in the documentation
  // for FB.getLoginStatus().
    if (response.status === 'connected') {
    // Logged into your app and Facebook.
      this.props.logined(response.authResponse.accessToken);
      // this.testAPI();
    } else if (response.status === 'not_authorized') {
    // The person is logged into Facebook, but not your app.
    //document.getElementById('status').innerHTML = 'Please log ' +
    //  'into this app.';
    } else {
    // The person is not logged into Facebook, so we're not sure if
    // they are logged into this app or not.
    //  document.getElementById('status').innerHTML = 'Please log ' +
    //'into Facebook.';
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
      <div>
        <a href="#" onClick={this.handleClick}>Login</a>
      </div>
    );
  }
});

React.render(
  <div>
    <LoginPage />
  </div>,
  document.getElementById('all')
);
