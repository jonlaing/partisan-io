import React from 'react';

export default React.createClass({
  getInitialState() {
    return {};
  },

  componentDidMount() {
  },

  render() {
    return (
      <div className="login">
        <div className="form-input">
          <input type="text" rel="username" />
        </div>
        <div className="form-input">
          <input type="password" rel="password" />
        </div>
        <button onClick={this.handleLogin}>Login</button>
      </div>
    );
  }
});
