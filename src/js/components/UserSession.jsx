import React from 'react';
import Icon from 'react-fontawesome';

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

  shouldComponentUpdate() {
    return false;
  },

  render() {
    return (
      <div className={"usersession " + this.props.className} >
        <a href="javascript:void(0)" className="usersession-name" onClick={this.handleDropdown} >
          <img src={this.props.avatar} height="27" width="27" />
          @{this.props.username} <Icon name="chevron-down" />
        </a>
        <ul className="usersession-dropdown" ref="dropdown">
          <li>
            <a href="javascript:void(0)" onClick={this.handleLogout}>Logout</a>
          </li>
        </ul>
      </div>
    );
  }
});
