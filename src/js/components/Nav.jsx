import React from 'react';

import Notifications from './Notifications.jsx';

export default React.createClass({
  render() {
    return (
      <div className="nav">
        <ul>
          <li>
            <a href="/feed/">Feed</a>
          </li>
          <li>
            <a href="/matches/">Matches</a>
          </li>
          <li>
            <a href="/questions/">Questions</a>
          </li>
          <li>
            <a href="/friends/">Friends</a>
          </li>
          <li>
            <Notifications />
          </li>
        </ul>
      </div>
    );
  }
});
