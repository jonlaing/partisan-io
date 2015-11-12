import React from 'react';
import moment from 'moment';

import Icon from 'react-fontawesome';

import NotificationActionCreator from '../actions/NotificationActionCreator';
import NotificationStore from '../stores/NotificationStore';

import Breakout from './Breakout.jsx';

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
    NotificationActionCreator.getNotificationList();
    this.setState({showList: !this.state.showList});
  },

  render() {
    var notifCount, fullList;
    if (this.state.count > 0) {
      notifCount = (
        <span className="notification-number">{this.state.count}</span>
      );
    } else {
      notifCount = "";
    }

    var list = this.state.notifications.map((notif) => {
      return (
        <li key={notif.notification.id}>
          <div className="notification-unseen">
            {this._unseenMarker(notif.notification)}
          </div>
          {this._notifTemplate(notif)}
        </li>
      );
    });

    if(this.state.showList === true) {
      fullList = (
        <Breakout className="notification-list">
          <ul>
            {list}
          </ul>
        </Breakout>
      );
    }


    return (
      <div className="notification-counter">
        <a href="javascript:void(0)" onClick={this.handleClick}>Notifications{notifCount}</a>
        {fullList}
      </div>
    );
  },

  _notifTemplate(notif) {
    switch(notif.notification.record_type) {
      case "comment":
        return this._commentTemplate(notif);
      case "like":
        return this._likeTemplate(notif);
      case "friendship":
        return this._friendTemplate(notif);
      case "user_tag":
        return this._userTagTemplate(notif);
      default:
        break;
    }
  },

  _commentTemplate(notif) {
    let username = notif.user.username;
    let route = "/comments/" + notif.notification.record_id;

    return (
      <span className={notif.notification.seen ? "seen" : "unseen" }>
        {this._avatarTemplate(notif)}
        <div>
          <a href={route}>
            @{username} commented on your post.
          </a>
          <small>{moment(notif.notification.created_at).fromNow()}</small>
        </div>
      </span>
    );
  },

  _likeTemplate(notif) {
    let username = notif.user.username;
    let route = "/likes/" + notif.notification.record_id;

    return (
      <span className={notif.notification.seen ? "seen" : "unseen" }>
        {this._avatarTemplate(notif)}
        <div>
          <a href={route}>
            @{username} liked your thing.
          </a>
          <small>{moment(notif.notification.created_at).fromNow()}</small>
        </div>
      </span>
    );
  },

  _friendTemplate(notif) {
    let username = notif.user.username;
    let route = "/profiles/" + username;

    if (notif.record.confirmed === false) {
      return (
        <span className={notif.notification.seen ? "seen" : "unseen" }>
          {this._avatarTemplate(notif)}
          <div>
            <a href={route}>
              @{username} sent you a friend request.
            </a>
            <small>{moment(notif.notification.created_at).fromNow()}</small>
          </div>
        </span>
      );
    } else {
      return (
        <span className={notif.notification.seen ? "seen" : "unseen" }>
          {this._avatarTemplate(notif)}
          <div>
            <a href={route}>
              @{username} confirmed your friendship.
            </a>
            <small>{moment(notif.notification.created_at).fromNow()}</small>
          </div>
        </span>
      );
    }
  },

  _userTagTemplate(notif) {
    console.log(notif);
    let username = notif.user.username;
    let route = "/posts/" + notif.record.record_id;

    return (
      <span className={notif.notification.seen ? "seen" : "unseen" }>
        {this._avatarTemplate(notif)}
        <div>
          <a href={route}>
            @{username} mentioned you in a post.
          </a>
          <small>{moment(notif.notification.created_at).fromNow()}</small>
        </div>
      </span>
    );
  },

  _avatarTemplate(notif) {
    return (
      <div className="notification-avatar">
        <img className="user-avatar" src={notif.user.avatar_thumbnail_url} />
      </div>
    );
  },

  _unseenMarker(notif) {
    if(notif.seen === false) {
      return <Icon name="circle" />;
    }

    return (<span>&nbsp;</span>);
  },

  _onChange() {
    let state = NotificationStore.getAll();
    this.setState(state);
  }

});
