var PageAll = React.createClass({
  getInitialState: function() {
    return {
      login: false
    };
  },
  logined: function() {
    this.setState({
      login: true
    });
  },
  render: function() {
    if (this.state.login == true) {
      return (
        <App/>
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
