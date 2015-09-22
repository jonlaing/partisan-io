import React from 'react';

import NotificationActionCreator from '../actions/NotificationActionCreator';
import NotificationStore from '../stores/NotificationStore';

export default React.createClass({
  getInitialState() {
    return {count: 0, notifications: [], showList: false};
  },

  componentDidMount() {
    NotificationStore.addChangeListener(this._onChange);
    NotificationActionCreator.getNotificationCount();
  },

  componentWillUnmount() {
    NotificationStore.removeChangeListener(this._onChange);
  },

  handleClick() {
    if(this.state.notifications.length === 0 && this.state.count > 0) {
      NotificationActionCreator.getNotificationList();
    }

    this.setState({showList: !this.state.showList});
  },

  render() {
    var notifCount, fullList;
    if (this.state.count > 0) {
      notifCount = (
        <span className="notification-number">({this.state.count})</span>
      );
    } else {
      notifCount = "";
    }

    var list = this.state.notifications.map((notif) => {
      return (
        <li key={notif.notification.id}>
          {this._notifTemplate(notif)}
        </li>
      );
    });

    if(this.state.showList === true) {
      fullList = (
        <ul className="notification-list">
          {list}
        </ul>
      );
    }


    return (
      <div className="notification-counter">
        <a href="javascript:void(0)" onClick={this.handleClick}>Notifications&nbsp;{notifCount}</a>
        {fullList}
      </div>
    );
  },

  _notifTemplate(notif) {
    switch(notif.notification.record_type) {
      case "comment":
      case "comments":
        return this._commentTemplate(notif);
      case "like":
      case "likes":
        return this._likeTemplate(notif);
      default:
        break;
    }
  },

  _commentTemplate(notif) {
    let username = notif.user.username;
    let route = "/comments/" + notif.notification.record_id;

    return (
      <span className={notif.notification.seen ? "seen" : "unseen" }>
        <a href={route}>
          @{username} commented on your post.
        </a>
        <small>{notif.notification.created_at}</small>
      </span>
    );
  },

  _likeTemplate(notif) {
    let username = notif.user.username;
    let route = "/likes/" + notif.notification.record_id;

    return (
      <span className={notif.notification.seen ? "seen" : "unseen" }>
        <a href={route}>
          @{username} liked your thing.
        </a>
        <small>{notif.notification.created_at}</small>
      </span>
    );
  },

  _onChange() {
    let state = NotificationStore.getAll();
    this.setState(state);
  }

});
