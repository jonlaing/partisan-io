import React from 'react';
import moment from 'moment';

import formatter from '../utils/formatter';

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
    var image;

    if(this.props.data.image_attachment.id !== 0) {
      image = (
        <img src={this.props.data.image_attachment.image_url} width="30%" />
      );
    } else {
      image = "";
    }

    return (
      <div className="comment">
        <div className="comment-avatar">
          <img className="user-avatar" src={formatter.avatarUrl(this.props.data.user.avatar_thumbnail_url)} />
        </div>
        <div>
          <div className="comment-author">
            <a href={"/profiles/" + this.props.data.user.id}>@{this.props.data.user.username}</a>
          </div>
          <div className="comment-body" dangerouslySetInnerHTML={formatter.comment(this.props.data.comment.body)} />
          <div className="comment-attachment">
            {image}
          </div>
        </div>
        <div className="comment-actions">
          <div className="right comment-meta">{moment(this.props.data.comment.created_at).fromNow()}</div>
          <Likes onClick={this.handleLike} count={this.props.data.like_count} liked={this.props.data.liked} />
        </div>
      </div>
    );
  }
});
