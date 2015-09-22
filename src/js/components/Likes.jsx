import React from 'react';

export default React.createClass({
  handleLike() {
    this.props.onClick();
  },

  render() {
    return (
      <div>
        <button className={"like" + (this.props.liked ? " active" : "")} onClick={this.handleLike}>
          <i className="fi-like"></i>
          {this.props.count}
        </button>
      </div>
    );
  }
});
