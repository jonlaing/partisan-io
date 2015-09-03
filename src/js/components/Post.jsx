import React from 'react';

import PostActions from './PostActions.jsx'; // TODO: Rename this dumb shit

export default React.createClass({

  render() {
    return (
      <div className="post">
        <div className="post-header">
          <div className="post-avatar">
            <img className="user-avatar" src="" />
          </div>
          <div className="post-user">
            <h4 className="post-username">@{this.props.data.user.username}</h4>
            <span className="post-timestamp">{this.props.data.post.created_at}</span>
          </div>
        </div>
        <div className="post-content">
          {this.props.data.post.body}
        </div>
        <PostActions id={this.props.data.post.id} likeCount={this.props.data.like_count} dislikeCount={this.props.data.dislike_count} commentCount={this.props.data.comment_count} />
      </div>
    );
  }
});
