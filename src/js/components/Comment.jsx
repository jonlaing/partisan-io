import React from 'react';
import moment from 'moment';

import LikeActionCreator from '../actions/LikeActionCreator';
import Likes from './Likes.jsx';

export default React.createClass({
  getInitialState() {
    return {};
  },

  handleLike() {
    LikeActionCreator.like("comment", this.props.data.comment.id);
  },

  render() {
    return (
      <div className="comment">
        <div className="comment-author">
          <a href={"/profiles/" + this.props.data.user.id}>@{this.props.data.user.username}</a>
        </div>
        <div className="comment-body">
          {this.props.data.comment.body}
        </div>
        <div>
          <div className="right comment-meta">{moment(this.props.data.comment.created_at).fromNow()}</div>
          <Likes onClick={this.handleLike} count={this.props.data.like_count} liked={this.props.data.liked} />
        </div>
      </div>
    );
  }
});
