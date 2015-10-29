import React from 'react';
import Icon from 'react-fontawesome';

import LoginActionCreator from '../actions/LoginActionCreator';

import Breakout from './Breakout.jsx';

export default React.createClass({
  getInitialState() {
    return { show: false };
  },

  handleLogout() {
    LoginActionCreator.logout();
    window.location.href = "/login";
  },

  handleDropdown() {
    this.setState({show: !this.state.show});
  },

  render() {
    return (
      <div className={"usersession " + this.props.className} >
        <a href="javascript:void(0)" className="usersession-name" onClick={this.handleDropdown} >
          <img src={this.props.avatar} height="27" width="27" />
          @{this.props.username} <Icon name="chevron-down" />
        </a>
        <Breakout className="usersession-dropdown" show={this.state.show}>
          <ul>
            <li>
              <a href="javascript:void(0)" onClick={this.handleLogout}>Logout</a>
            </li>
          </ul>
        </Breakout>
      </div>
    );
  }
});
