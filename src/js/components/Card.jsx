import React from 'react';

export default React.createClass({
  render() {
    return (
      <div className="card">
        {this.props.children}
      </div>
    );
  }
});
