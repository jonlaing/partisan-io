import React from 'react';

import NotificationActionCreator from '../actions/NotificationActionCreator';
import NotificationStore from '../stores/NotificationStore';

// Polling for changes. This seemed like the best place for it, since I'm thinking
// of staying on the page as a user input, albeit passive input. This might not
// be the right place though. Will revisit.

export default React.createClass({
  getInitialState() {
    return {count: 0, notifications: []};
  },

  componentDidMount() {
    NotificationStore.addChangeListener(this._onChange);
    NotificationActionCreator.getNotificationCount();
  },

  componentWillUnmount() {
    NotificationStore.removeChangeListener(this._onChange);
  },

  render() {
    var list = this.state.notifications.map((notif) => {
      return (
        <li key={notif.notification.id}>
          {this._notifTemplate(notif)}
        </li>
      );
    });

    return (
      <div className="notification-counter">
        Notifications
        <span className="notification-number">{this.state.count}</span>
        <ul className="notification-list">
          {list}
        </ul>
      </div>
    );
  },

  _notifTemplate(notif) {
    switch(notif.notification.record_type) {
      case "comment":
      case "comments":
        return this._commentTemplate(notif);
      default:
        break;
    }
  },

  _commentTemplate(notif) {
    let username = notif.user.username;

    return (
      <span className={notif.notification.seen ? "seen" : "unseen" }>
        @{username} commented on your post.
        <small>{notif.notification.created_at}</small>
      </span>
    );
  },

  _onChange() {
    let state = NotificationStore.getCount();
    this.setState(state);
  }

});
