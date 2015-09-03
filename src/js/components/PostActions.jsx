import React from 'react';

export default React.createClass({
  getInitialState() {
    return {};
  },

  componentDidMount() {
  },

  render() {
    return (
      <div className="post-actions">
        <button className="like" onClick={this.props.onLike}>
          {this.props.likeCount}
        </button>
        <button className="dislike" onClick={this.props.onDislike}>
          {this.props.dislikeCount}
        </button>
        <button className="comment" onClick={this.props.onComment}>
          {this.props.commentCount}
        </button>
      </div>
    );
  }
});
