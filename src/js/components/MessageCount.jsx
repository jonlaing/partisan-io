import React from 'react';

import MessageActionCreator from '../actions/MessageActionCreator';
import MessageStore from '../stores/MessageStore';

export default React.createClass({
  getInitialState() {
    return {count: 0};
  },

  componentDidMount() {
    MessageStore.addChangeListener(this._onChange);
    MessageActionCreator.getMessageCount();
  },

  componentWillUnmount() {
    MessageStore.removeChangeListener(this._onChange);
  },

  render() {
    var msgCount;
    if (this.state.count > 0) {
      msgCount = (
        <span className="message-number">{this.state.count}</span>
      );
    } else {
      msgCount = "";
    }

    return (
      <div className="message-counter">
        <a href="/messages">Messages{msgCount}</a>
      </div>
    );
  },

  _onChange() {
    let count = MessageStore.getCount();
    this.setState({count: count});
  }

});
