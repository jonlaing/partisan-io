import React from 'react';
import Icon from 'react-fontawesome';

export default React.createClass({
  handleLike() {
    this.props.onClick();
  },

  render() {
    return (
      <div style={{display: 'inline-block'}}>
        <button className={"like" + (this.props.liked ? " active" : "")} onClick={this.handleLike}>
          <Icon name="star" />
          Favorite&nbsp;
          <strong>({this.props.count})</strong>
        </button>
      </div>
    );
  }
});
