import React from 'react';
import Icon from 'react-fontawesome';

export default React.createClass({
  render() {
    return (
      <div className={this.props.className}>
        <button onClick={this.props.onClick} >
          <Icon name="comments-o" />
          Comment&nbsp;
          <strong>({this.props.count})</strong>
        </button>
      </div>
    );
  }
});
