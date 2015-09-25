import React from 'react';
import LoginActionCreator from '../actions/LoginActionCreator';

export default React.createClass({
  handleLogout() {
    LoginActionCreator.logout();
    window.location.href = "/login.html";
  },

  handleDropdown() {
    let dropdown = $(React.findDOMNode(this.refs.dropdown));
    dropdown.toggleClass('show');
  },

  render() {
    return (
      <div className="usersession">
        <a href="javascript:void(0)" className="usersession-name" onClick={this.handleDropdown} >@{this.props.username}</a>
        <ul className="usersession-dropdown" ref="dropdown">
          <li>
            <a href="javascript:void(0)" onClick={this.handleLogout}>Logout</a>
          </li>
        </ul>
      </div>
    );
  }
});
