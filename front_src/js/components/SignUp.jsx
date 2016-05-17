import React from 'react';
import ReactDOM from 'react-dom';
import SignUpActionCreator from '../actions/SignUpActionCreator';
import SignUpStore from '../stores/SignUpStore';

export default React.createClass({
  getInitialState() {
    return { errors: [], success: false, userUnique: 0 };
  },

  handleSubmit() {
    let email = $(ReactDOM.findDOMNode(this.refs.email)).val();

    let user = { email: email };

    SignUpActionCreator.signUp(user);
  },

  handleUsernameChange(e) {
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
      uniquenessMarker = <span className="signup-uniqueness success"><i className="fi-check"></i></span>;
    } else if (this.state.userUnique === 2) {
      uniquenessMarker = <span className="signup-uniqueness alert"><i className="fi-x"></i></span>;
    } else {
      uniquenessMarker = '';
    }

    return (
      <div className="signup">
        <div className={"form-input" + this._hasError("email")}>
          <input type="text" placeholder="you@email.com" ref="email" />
          {this._error("email")}
        </div>
        <div className={"form-input" + this._hasError("username")}>
          <div className="prefixed">
            <input type="text" placeholder="Username" ref="username" onChange={this.handleUsernameChange} />
            <span className="prefix">@</span>
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

        <div className="actions">
          <div className="right">
            <a href="/login">Login</a>
          </div>
          <button onClick={this.handleSubmit}>Sign Up</button>
        </div>
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
