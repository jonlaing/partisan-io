import React from 'react';

import formatter from '../utils/formatter';

import Friend from './Friend.jsx';
import UserSession from './UserSession.jsx';
import Notifications from './Notifications.jsx';

export default React.createClass({
  getInitialState() {
    return {};
  },

  componentDidMount() {
  },

  render() {
    return (
      <div className="profile">
        <header>
          <UserSession username={this.props.user.username} />
          <Notifications />
        </header>

        <div className="profile-container">
          <div className="profile-user">
            <h1>@{this.props.user.username}</h1>
            <div className="right">
              <Friend id={this.props.user.id} />
            </div>
          </div>
          <div className="profile-match">
            {this.props.match}%<span>Match</span>
          </div>
          <div className="profile-info">
            <div className="profile-info-location">{this.props.user.location}</div>
            <div className="profile-info-gender">{this.props.user.gender}</div>
          </div>
          <div className="profile-summary">
            {formatter.userSummary(this.props.profile.summary)}
          </div>
        </div>
      </div>
    );
  }
});
