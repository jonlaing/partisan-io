import React from 'react';
import moment from 'moment';

import PostLikes from './PostLikes.jsx';

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
              <h4 className="post-username">
                <a href={"/profiles/" + this.props.data.user.id}>@{this.props.data.user.username}</a>
              </h4>
              <span className="post-timestamp">{moment(this.props.data.post.created_at).fromNow()}</span>
            </div>
          </div>
          <div className="post-body">
            {this.props.data.post.body}
          </div>
        </div>
        <div className="post-actions">
          <PostLikes id={this.props.data.post.id} />
        </div>
      </div>
    );
  }
});
