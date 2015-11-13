import React from 'react';

import Message from './Message.jsx';

export default React.createClass({
  getInitialState() {
    return {};
  },

  componentDidMount() {
    this._scroll();
  },

  shouldComponentUpdate(nextProps) {
    return nextProps.messages.length !== this.props.messages.length;
  },

  componentDidUpdate() {
    this._scroll();
  },

  render() {
    var messages = this.props.messages.map((msg) => {
      return <Message key={msg.thread_id + '-' + msg.id} message={msg} thisUser={this.props.userID === parseInt(msg.user_id)} />;
    });

    return (
      <div className="message-list" ref="list">
        <ul>
          {messages}
        </ul>
      </div>
    );
  },

  _scroll() {
    let list = $(React.findDOMNode(this.refs.list));
    let height = list.find("ul").outerHeight();
    list.scrollTop(height);
  }
});
