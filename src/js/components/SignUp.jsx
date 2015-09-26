import React from 'react';
import SignUpActionCreator from '../actions/SignUpActionCreator';
import SignUpStore from '../stores/SignUpStore';

export default React.createClass({
  getInitialState() {
    return { errors: [], success: false, userUnique: 0 };
  },

  handleSubmit() {
    let email = $(React.findDOMNode(this.refs.email)).val();
    let username = $(React.findDOMNode(this.refs.username)).val();
    let postalCode = $(React.findDOMNode(this.refs.postalCode)).val();
    let password = $(React.findDOMNode(this.refs.password)).val();
    let passwordConfirm = $(React.findDOMNode(this.refs.passwordConfirm)).val();

    let user = {
      email: email,
      username: username,
      postalCode: postalCode,
      password: password,
      passwordConfirm: passwordConfirm
    };

    SignUpActionCreator.signUp(user);
  },

  handleUsernameChange(e) {
    console.log(e.target.value);
    SignUpActionCreator.checkUnique(e.target.value);
  },

  componentDidMount() {
    SignUpStore.addChangeListener(this._onChange);
  },

  componentWillUnmount() {
    SignUpStore.removeChangeListener(this._onChange);
  },

  componentDidUpdate() {
    if(this.state.success === true) {
      window.location.href = "/questions";
    }
  },

  render() {
    var uniquenessMarker;

    if(this.state.userUnique === 1) {
      uniquenessMarker = <span className="label success"><i className="fi-check"></i></span>;
    } else if (this.state.userUnique === 2) {
      uniquenessMarker = <span className="label alert"><i className="fi-x"></i></span>;
    }

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
            <input type="text" placeholder="Username" ref="username" onChange={this.handleUsernameChange} />
            {uniquenessMarker}
          </div>
          {this._error("username")}
        </div>
        <div className="form-input">
          <input type="text" placeholder="Postal Code ex: 11211" ref="postalCode" />
        </div>
        <div className="form-input">
          <input type="password" placeholder="Password" ref="password" />
        </div>
        <div className={"form-input" + this._hasError("password_confirm")}>
          <input type="password" placeholder="Password Confirm" ref="passwordConfirm" />
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
