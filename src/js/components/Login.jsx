/*global $ */
import React from 'react/addons';
import LoginActionCreator from '../actions/LoginActionCreator';
import LoginStore from '../stores/LoginStore';

var ReactCSSTransitionGroup = React.addons.CSSTransitionGroup;

function getStateFromStore() {
  return LoginStore.getLoginState();
}

export default React.createClass({
  getInitialState() {
    return {error: "", success: false };
  },

  handleLogin() {
    let email = $(React.findDOMNode(this.refs.email)).val();
    let password = $(React.findDOMNode(this.refs.password)).val();
    LoginActionCreator.login(email, password);
  },

  componentDidMount() {
    LoginStore.addChangeListener(this._onChange);
  },

  componentWillUnmount() {
    LoginStore.removeChangeListener(this._onChange);
  },

  componentDidUpdate() {
    if(this.state.success === true) {
      window.location.href = "/feed";
    }
  },

  render() {
    var error;

    if(this.state.error !== "") {
      error = (<div className="error">{this.state.error}</div>);
    } else {
      error = "";
    }

    return (
      <div className="login">
        <h4>Login to Partisan.IO</h4>
        <ReactCSSTransitionGroup transitionName="login-error">
          {error}
        </ReactCSSTransitionGroup>
        <div className="form-input">
          <input type="text" placeholder="you@email.com" ref="email" />
        </div>
        <div className="form-input">
          <input type="password" placeholder="Password" ref="password" />
        </div>
        <div className="right">
          <a href="/sign-up.html">Sign Up</a>
        </div>
        <button onClick={this.handleLogin}>Login</button>
      </div>
    );
  },
  _onChange() {
    this.setState(getStateFromStore());
  }
});
