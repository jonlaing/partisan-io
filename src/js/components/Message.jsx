import React from 'react';

import moment from 'moment';

import Breakout from './Breakout.jsx';

export default React.createClass({
  getInitialState() {
    return {};
  },

  handleAvatarClick(username) {
    if(this.props.thisUser === true) {
      return function() {};
    }

    return function() {
      window.location.href = "/profiles/" + username;
    };
  },

  componentDidMount() {
  },

  render() {
    var className = "message";
    if(this.props.thisUser === true) {
      className += " message-mine";
    }

    // make sure newlines are respected
    var text = this.props.message.body.split("\n").map((line, i) => <span key={i}>{line}<br/></span>);

    return (
      <li className={className}>
        <Breakout>
          <div className="message-avatar">
            <img src={this.props.message.user.avatar_thumbnail_url} className="user-avatar" onClick={this.handleAvatarClick(this.props.message.user.username)} />
          </div>
          <div className="message-body">
            <div className="message-timestamp">{moment(this.props.message.created_at).fromNow()}</div>
            {text}
          </div>
        </Breakout>
      </li>
    );
  }
});
