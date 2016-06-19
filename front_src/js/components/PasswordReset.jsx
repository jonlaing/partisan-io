import React from 'react';
import ReactDOM from 'react-dom';
import fetch from 'whatwg-fetch';

export default React.createClass({
  getInitialState() {
    return { errors: [], success: false };
  },

  handleSubmit() {
    let email = $(ReactDOM.findDOMNode(this.refs.email)).val();
    let pw = $(ReactDOM.findDOMNode(this.refs.password)).val();
    let pwC = $(ReactDOM.findDOMNode(this.refs.passwordConfirm)).val();

    fetch(`api/v2/password_reset/${this.props.resetID}`, {
      method: 'POST',
      body: JSON.stringify({
        email: email,
        password: pw,
        password_confirm: pwC
      })
    })
    .then(() => this.setState({succes: true}))
    .catch(err => console.log(err));
  },

  render() {
    return (
      <div className="passwordreset">
        <div className="form-input">
          <input type="text" placeholder="you@email.com" ref="email" />
        </div>
        <div className="form-input">
          <input type="password" placeholder="Password" ref="password" />
        </div>
        <div className={"form-input"}>
          <input type="password" placeholder="Password Confirm" ref="passwordConfirm" />
        </div>

        <div className="actions">
          <button onClick={this.handleSubmit}>Reset Password</button>
        </div>
      </div>
    );
  }
});
