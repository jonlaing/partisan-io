import React from 'react';
import moment from 'moment';

import PostActions from './PostActions.jsx'; // TODO: Rename this dumb shit

export default React.createClass({

  render() {
    return (
      <div className="post">
        <div className="card-body">
          <div className="post-header">
            <div className="post-avatar">
              <img className="user-avatar" src="" />
            </div>
            <div className="post-user">
              <h4 className="post-username">@{this.props.data.user.username}</h4>
              <span className="post-timestamp">{moment(this.props.data.post.created_at).fromNow()}</span>
            </div>
          </div>
          <div className="post-body">
            {this.props.data.post.body}
          </div>
        </div>
        <PostActions id={this.props.data.post.id} likeCount={this.props.data.like_count} dislikeCount={this.props.data.dislike_count} commentCount={this.props.data.comment_count} />
      </div>
    );
  }
});
