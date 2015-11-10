import React from 'react';

import ThreadActionCreator from '../actions/ThreadActionCreator';
import MessageActionCreator from '../actions/MessageActionCreator';
import ThreadStore from '../stores/ThreadStore';
import MessageStore from '../stores/MessageStore';

import Nav from './Nav.jsx';
import UserSession from './UserSession.jsx';
import ThreadList from './ThreadList.jsx';
import MessageList from './MessageList.jsx';
import MessageComposer from './MessageComposer.jsx';

export default React.createClass({
  getInitialState() {
    return {threads: [], currentThread: null, messages: []};
  },

  componentDidMount() {
    ThreadStore.addChangeListener(this._onThreadsChange);
    MessageStore.addChangeListener(this._onMessagesChange);
    ThreadActionCreator.getThreads();
    MessageActionCreator.getMessages(1);
    MessageActionCreator.messageSocket(1);
  },

  componentWillUnmount() {
    ThreadStore.removeChangeListener(this._onThreadsChange);
    MessageStore.removeChangeListener(this._onMessagesChange);
  },

  render() {
    return (
      <div className="messages">
        <header className="header">
          <UserSession className="right" username={this.props.data.user.username} avatar={this.props.data.user.avatar_thumbnail_url} />
          <img src="images/logo.svg" className="logo" />
          <Nav currentPage="messages" />
        </header>

        <div className="dashboard">
          <aside>
            <ThreadList threads={this.state.threads} />
          </aside>
          <article className="messages-container">
            <MessageList messages={this.state.messages} userID={this.props.data.user.id} />
            <MessageComposer thread={this.state.currentThread}/>
          </article>
        </div>
      </div>
    );
  },

  _onThreadsChange() {
    this.setState({
      threads: ThreadStore.getList()
    });
  },

  _onMessagesChange() {
    this.setState({
      messages: MessageStore.getMessages()
    });
  }
});
