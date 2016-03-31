import React from 'react';

import formatter from '../utils/formatter';

export default React.createClass({
  getInitialState() {
    return {};
  },

  componentDidMount() {
  },

  render() {
    return (
      <div className="match">
        <div className="match-avatar">
          <img src={formatter.avatarUrl(this.props.user.avatar_thumbnail_url)} className="user-avatar" />
        </div>
        <div>
          <div className="match-user">
            <a href={"profiles/" + this.props.user.username}>@{this.props.user.username}</a>
          </div>
          <div className="match-info">{formatter.age(this.props.user.birthdate)}&nbsp;-&nbsp;{this.props.user.gender}</div>
          <div className="match-location">{formatter.cityState(this.props.user.location)}</div>
          <div className="match-match">{formatter.match(this.props.match)}</div>
        </div>
      </div>
    );
  }
});
