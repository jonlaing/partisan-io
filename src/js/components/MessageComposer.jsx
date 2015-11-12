import React from 'react';

import MessageActionCreator from '../actions/MessageActionCreator';

const _ENTER = 13;

export default React.createClass({
  getInitialState() {
    return {};
  },

  handleKeyDown(e) {
    if(this.props.thread === 0) {
      e.target.value = "";
      return;
    }

    if(e.keyCode === _ENTER && !e.shiftKey) {
      e.preventDefault();
      MessageActionCreator.sendMessage(this.props.thread, e.target.value);
      e.target.value = "";
    }
  },

  handleButtonClick() {
    var message = React.findDOMNode(this.refs.message);

    if(this.props.thread === 0) {
      message.value = "";
      return;
    }

    MessageActionCreator.sendMessage(this.props.thread, message.value);
    message.value = "";
  },

  componentDidMount() {
  },

  render() {
    return (
      <div className="message-composer">
        <textarea type="text" placeholder="Type your message here" onKeyDown={this.handleKeyDown} ref="message"/>
        <button onClick={this.handleButtonClick}>Send</button>
      </div>
    );
  }
});
