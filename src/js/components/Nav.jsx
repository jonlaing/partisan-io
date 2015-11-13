import React from 'react';

import Notifications from './Notifications.jsx';
import MessageCount from './MessageCount.jsx';

export default React.createClass({
  render() {
    return (
      <div className="nav">
        <ul>
          <li>
            <a href="/feed/" className={ this.props.currentPage === "feed" ? "active" : "" }>Feed</a>
          </li>
          <li>
            <a href="/matches/" className={ this.props.currentPage === "matches" ? "active" : "" }>Matches</a>
          </li>
          <li>
            <a href="/questions/" className={ this.props.currentPage === "questions" ? "active" : "" }>Questions</a>
          </li>
          <li>
            <a href="/friends/" className={ this.props.currentPage === "friends" ? "active" : "" }>Friends</a>
          </li>
          <li>
            <Notifications />
          </li>
          <li>
            <MessageCount />
          </li>
        </ul>
      </div>
    );
  }
});
