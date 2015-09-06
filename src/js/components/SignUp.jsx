/*global $ */
import React from 'react';
import SignUpActionCreator from '../actions/SignUpActionCreator';
import SignUpStore from '../stores/SignUpStore';

export default React.createClass({
  getInitialState() {
    return { errors: [], success: false };
  },

  handleSubmit() {
    let email = $(React.findDOMNode(this.refs.email)).val();
    let username = $(React.findDOMNode(this.refs.username)).val();
    let fullName = $(React.findDOMNode(this.refs.full_name)).val();
    let password = $(React.findDOMNode(this.refs.password)).val();
    let passwordConfirm = $(React.findDOMNode(this.refs.password_confirm)).val();

    let user = {
      email: email,
      username: username,
      full_name: fullName,
      password: password,
      password_confirm: passwordConfirm
    };

    SignUpActionCreator.signUp(user);
  },

  componentDidMount() {
    SignUpStore.addChangeListener(this._onChange);
  },

  componentWillUnmount() {
    SignUpStore.removeChangeListener(this._onChange);
  },

  componentDidUpdate() {
    if(this.state.success === true) {
      window.location.href = "/feed.html";
    }
  },

  render() {
    return (
      <div className="signup">
        <h4>Sign Up for Partisan.IO</h4>
        <div className={"form-input" + this._hasError("email")}>
          <input type="text" placeholder="you@email.com" ref="email" />
          {this._error("email")}
        </div>
        <div className={"form-input row collapse" + this._hasError("username")}>
          <div className="large-1 columns">
            <span className="prefix">@</span>
          </div>
          <div className="large-11 columns">
            <input type="text" placeholder="Username" ref="username" />
          </div>
          {this._error("username")}
        </div>
        <div className="form-input">
          <input type="text" placeholder="Full Name" ref="full_name" />
        </div>
        <div className="form-input">
          <input type="password" placeholder="Password" ref="password" />
        </div>
        <div className={"form-input" + this._hasError("password_confirm")}>
          <input type="password" placeholder="Password Confirm" ref="password_confirm" />
          {this._error("password_confirm")}
        </div>
        <button onClick={this.handleSubmit}>Sign Up</button>
      </div>
    );
  },

  _onChange() {
    this.setState(SignUpStore.getState());
  },

  _hasError(field) {
    if(this.state.errors[field]) {
      return " error";
    }
    return "";
  },

  _error(field) {
    let err = this.state.errors[field];
    if(err) {
      return (
        <small className="error">{err}</small>
      );
    }
  }
});
