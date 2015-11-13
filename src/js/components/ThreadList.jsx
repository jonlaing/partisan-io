import React from 'react';

import ThreadActionCreator from '../actions/ThreadActionCreator';

export default React.createClass({
  getInitialState() {
    return {filter: ""};
  },

  handleUsernameChange(e) {
    this.setState({filter: e.target.value});
  },

  handleThreadSwitch(threadID) {
    return function() {
      ThreadActionCreator.switchThreads(threadID);
    };
  },

  handleThreadCreate(friendID) {
    return function() {
      ThreadActionCreator.createThread(friendID);
    };
  },

  shouldComponentUpdate(nextProps, nextState) {
    return nextProps.threads !== this.props.threads ||
      nextProps.currentThread !== this.props.currentThread ||
      nextState.filter !== this.state.filter;
  },

  render() {
    var threads = this.props.threads.map((thread) => {
      if(thread.thread_user.user.username.includes(this.state.filter)) {
        let t = thread.thread_user;
        let className = t.thread_id === this.props.currentThread ? "thread-selected" : "";
        className += thread.has_unread ? "thread-unread" : "";

        return (
          <li key={t.thread_id} className={className} onClick={this.handleThreadSwitch(t.thread_id)}>
            <div className="thread-avatar">
              <img src={t.user.avatar_thumbnail_url} className="user-avatar" />
            </div>
            <div>
              {t.user.username}
            </div>
          </li>
       );
      }
    });

    var inactive = this.props.inactive.map((user) => {
      if(user.username.includes(this.state.filter)) {
        return (
          <li key={user.id} onClick={this.handleThreadCreate(user.id)} >
            <div className="thread-avatar">
              <img src={user.avatar_thumbnail_url} className="user-avatar" />
            </div>
            <div>
              {user.username}
            </div>
          </li>
        );
      }
    });

    return (
      <div className="thread-list">
        <div className="thread-list-header">
          <input type="text" placeholder="Search for users" onChange={this.handleUsernameChange} />
        </div>
        <ul>
          {threads}
          {inactive}
        </ul>
      </div>
    );
  }
});
