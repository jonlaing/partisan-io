import React from 'react';
import LoginActionCreator from '../actions/LoginActionCreator';
import LoginStore from '../stores/LoginStore';

export default React.createClass({
  getInitialState() {
    return { user: { username: "" } };
  },

  handleLogout() {
    LoginActionCreator.logout();
    window.location.href = "/login.html";
  },

  handleDropdown() {
    let dropdown = $(React.findDOMNode(this.refs.dropdown));
    dropdown.toggleClass('show');
  },

  componentDidMount() {
    LoginStore.addChangeListener(this._onChange);
    LoginActionCreator.fetchCurrentUser();
  },

  componentWillUnmount() {
    LoginStore.removeChangeListener(this._onChange);
  },

  render() {
    return (
      <div className="usersession">
        <a href="javascript:void(0)" className="usersession-name" onClick={this.handleDropdown} >@{this.state.user.username}&nbsp;<i className="fi-widget"></i></a>
        <ul className="usersession-dropdown" ref="dropdown">
          <li>
            <a href="javascript:void(0)" onClick={this.handleLogout}>Logout</a>
          </li>
        </ul>
      </div>
    );
  },

  _onChange() {
    this.setState({user: LoginStore.getUser()});
  }
});
