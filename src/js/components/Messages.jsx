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
    return {threads: [], inactiveThreads: [], currentThread: 0, messages: []};
  },

  componentDidMount() {
    ThreadStore.addChangeListener(this._onThreadsChange);
    MessageStore.addChangeListener(this._onMessagesChange);
    ThreadActionCreator.getThreads();
  },

  componentWillUnmount() {
    ThreadStore.removeChangeListener(this._onThreadsChange);
    MessageStore.removeChangeListener(this._onMessagesChange);
  },

  componentDidUpdate(_, prevState) {
    if(prevState.currentThread !== this.state.currentThread) {
      MessageActionCreator.getMessages(this.state.currentThread);
      MessageActionCreator.messageSocket(this.state.currentThread);
      return;
    }

    // if no thread is set, set one
    if(this.state.currentThread === 0 && this.state.threads.length > 0) {
      this.setState({currentThread: this.state.threads[0].thread_user.thread_id});
    }
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
            <ThreadList threads={this.state.threads} inactive={this.state.inactiveThreads}/>
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
      threads: ThreadStore.getList(),
      inactiveThreads: ThreadStore.getInactive(),
      currentThread: ThreadStore.getCurrentThread()
    });
  },

  _onMessagesChange() {
    this.setState({
      messages: MessageStore.getMessages()
    });
  }
});
