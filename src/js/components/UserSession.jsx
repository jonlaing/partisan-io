import React from 'react';

export default React.createClass({
  getInitialState() {
    return {};
  },

  handleLogout() {
  },

  componentDidMount() {
  },

  render() {
    return (
      <div className="usersession">
        <span className="usersession-name">{this.props.user.username}</span>
        <a href="javascript:void(0)" onClick={this.handleLogout}>Logout</a>
      </div>
    );
  }
});
