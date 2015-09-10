import React from 'react';

export default React.createClass({
  getInitialState() {
    return {};
  },

  componentDidMount() {
  },

  render() {
    return (
      <div className="profile">
        <h1>@{this.props.user.username}</h1>
        <div className="row profile-match">
          <div className="large-6 columns profile-match-match">{this.props.match}%<span>Match</span></div>
          <div className="large-6 columns text-right profile-match-enemy">{this.props.enemy}%<span>Enemy</span></div>
        </div>
        <div className="row">
          <div className="large-6 columns">{this.props.user.location}</div>
          <div className="large-6 columns">{this.props.user.gender}</div>
        </div>
      </div>
    );
  }
});
