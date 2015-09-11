import React from 'react';

export default React.createClass({
  getInitialState() {
    return { likeCount: 0, dislikeCount: 0, commentCount: 0 };
  },

  componentDidMount() {
  },

  render() {
    return (
      <div className="post-actions clearfix">
        <button className="like" onClick={this.props.onLike}>
          <i className="fi-like"></i>
          {this.state.likeCount}
        </button>
        <button className="dislike" onClick={this.props.onDislike}>
          <i className="fi-dislike"></i>
          {this.state.dislikeCount}
        </button>
        <button className="comment right" onClick={this.props.onComment}>
          <i className="fi-comment"></i>
          {this.state.commentCount}
        </button>
      </div>
    );
  }
});
