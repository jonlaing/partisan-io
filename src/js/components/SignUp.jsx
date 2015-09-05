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
        <div className="form-input">
          <input type="text" placeholder="you@email.com" ref="email" />
        </div>
        <div className="form-input row collapse">
          <div className="large-1 columns">
            <span className="prefix">@</span>
          </div>
          <div className="large-11 columns">
            <input type="text" placeholder="Username" ref="username" />
          </div>
        </div>
        <div className="form-input">
          <input type="text" placeholder="Full Name" ref="full_name" />
        </div>
        <div className="form-input">
          <input type="password" placeholder="Password" ref="password" />
        </div>
        <div className="form-input">
          <input type="password" placeholder="Password Confirm" ref="password_confirm" />
        </div>
        <button onClick={this.handleSubmit}>Sign Up</button>
      </div>
    );
  },

  _onChange() {
    this.setState(SignUpStore.getState());
  }
});
