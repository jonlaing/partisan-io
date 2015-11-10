import React from 'react';

export default React.createClass({
  getInitialState() {
    return {};
  },

  render() {
    var threads = this.props.threads.map((thread) => {
      return <li key={thread.id}>{thread.thread_user.user.username}</li>;
    });

    return (
      <div className="thread-list">
        <div className="thread-list-header">
          <h2>Message Threads</h2>
          <button className="thread-list-add">Start Chat</button>
        </div>
        <ul>
          {threads}
        </ul>
      </div>
    );
  }
});
