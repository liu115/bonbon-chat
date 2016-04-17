MessageBalloon = React.createClass({
  getInitialState: function() {
    var find_url = this.findURL(); 
    if (find_url != null) {
      var url = find_url[0];
      var protocal = find_url[1];
      var remain = find_url[2];
      var httpRequest = new XMLHttpRequest();
      httpRequest.onreadystatechange = function() {
        if (httpRequest.readyState == 4) {
          if (httpRequest.status == 200) {
            var meta = JSON.parse(httpRequest.responseText);
            this.setState({meta: meta})
          }
        }
      }.bind(this)
      httpRequest.open('GET', window.location.origin + '/meta?protocal=' + protocal + "&url=" + remain);
      httpRequest.send(null);
    }
	  return {meta: null};
  },
  findURL: function() {
    var urlPattern = /(http|https):\/\/([\w-]+(\.[\w-]+)+([\w.,@?^=%&amp;:\/~+#-]*[\w@?^=%&amp;\/~+#-]))?/
    return this.props.msg.content.match(urlPattern);
  },
  showMsg: function() {
    var res = this.findURL();
    if (res != null) {
      var url = res[0];
      var start = res.index;
      var end = start + url.length;
      var left = this.props.msg.content.substring(0, start);
      var right = this.props.msg.content.substring(end);
      return (
        [left, <a target="_blank" href={url}>{url}</a>, right]
      )
    } else {
      return this.props.msg.content;
    }
  },
  showMeta: function() {
    if (this.state.meta != null) {
      return (
        [
          <hr />,
          this.state.meta.sitename ? (<div className={"meta-sitename"}>{this.state.meta.sitename}</div>) : "",
          this.state.meta.title ? (<div className={"meta-title"}>{this.state.meta.title}</div>) : "",
          this.state.meta.description ? (<div className={"meta-description"}>{this.state.meta.description}</div>) : "",
          this.state.meta.image ? (<img className={"meta-image"} src={this.state.meta.image} />) : ""
          ]
      )
    }
  },
  render: function() {
    return (
      <div className={"wrapper-message-" + this.props.msg.from}>
        <div className={"message-balloon message-" + this.props.msg.from}>
          {this.showMsg()}
          {this.showMeta()}
        </div>
      </div>
    )
  }
});

