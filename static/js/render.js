var PageAll = React.createClass({
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
        <App token={this.state.token}/>
      );
    }
    else {
      return (
        <div>
          <LoginPage logined={this.logined}/>
        </div>
      );
    }
  }
});

React.render(
  <div>
    <PageAll/>
  </div>,
  document.getElementById('all')
);
