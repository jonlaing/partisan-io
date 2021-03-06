/*global $ */
import React from 'react';
import ReactDOM from 'react-dom';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';

import LoginActionCreator from '../actions/LoginActionCreator';
import LoginStore from '../stores/LoginStore';

function getStateFromStore() {
  return LoginStore.getLoginState();
}

export default React.createClass({
  getInitialState() {
    return {error: "", success: false };
  },


  handleLogin(e) {
    e.preventDefault();
    let email = $(ReactDOM.findDOMNode(this.refs.email)).val();
    let password = $(ReactDOM.findDOMNode(this.refs.password)).val();
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
        <ReactCSSTransitionGroup transitionName="login-error" transitionEnterTimeout={1000} transitionLeaveTimeout={1000}>
          {error}
        </ReactCSSTransitionGroup>
        <form onSubmit={this.handleLogin}>
          <div>
            <input type="text" placeholder="you@email.com" ref="email" />
          </div>
          <div>
            <input type="password" placeholder="Password" ref="password" />
          </div>

          <div className="right">
            <a href="/signup">Sign Up</a>
          </div>
          <button onClick={this.handleLogin}>Login</button>
        </form>
      </div>
    );
  },
  _onChange() {
    this.setState(getStateFromStore());
  }
});
