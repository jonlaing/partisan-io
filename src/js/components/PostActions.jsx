import React from 'react';

export default React.createClass({
  getInitialState() {
    return {};
  },

  componentDidMount() {
  },

  render() {
    return (
      <div className="post-actions clearfix">
        <button className="like" onClick={this.props.onLike}>
          <i className="fi-like"></i>
          {this.props.likeCount}
        </button>
        <button className="dislike" onClick={this.props.onDislike}>
          <i className="fi-dislike"></i>
          {this.props.dislikeCount}
        </button>
        <button className="comment right" onClick={this.props.onComment}>
          <i className="fi-comment"></i>
          {this.props.commentCount}
        </button>
      </div>
    );
  }
});
