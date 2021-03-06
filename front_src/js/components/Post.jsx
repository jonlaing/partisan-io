import React from 'react';
import moment from 'moment';

import formatter from '../utils/formatter';

import Icon from 'react-fontawesome';

import LikeActionCreator from '../actions/LikeActionCreator';
import FlagActionCreator from '../actions/FlagActionCreator';
import LightboxActionCreator from '../actions/LightboxActionCreator';

import Likes from './Likes.jsx';
import CommentCounter from './CommentCounter.jsx';
import CommentList from './CommentList.jsx';

export default React.createClass({
  getInitialState() {
    return {showComments: this.props.defaultShowComments || false};
  },


  handleToggleComments() {
    let show = !this.state.showComments;
    this.setState({showComments: show});
  },

  handleLike() {
    LikeActionCreator.like("post", this.props.data.post.id);
  },

  handleFlag() {
    FlagActionCreator.beginReport(this.props.data.post.id, "post");
  },

  handleImageClick() {
    LightboxActionCreator.open(this.props.data.image_attachment.image_url);
  },

  render() {
    var comments, attachment;

    if(this.state.showComments === true) {
      comments = (
        <div>
          <CommentList id={this.props.data.post.id} show={this.state.showComments}/>
        </div>
      );
    }

    if(this.props.data.image_attachment.id > 0) {
      attachment = (
        <div className="post-attachment">
          <img src={this.props.data.image_attachment.image_url} width="100%" onClick={this.handleImageClick}/>
        </div>
      );
    }

    return (
      <div className="post">
        <div className="card-body">
          <div className="post-header">
            <div className="post-avatar">
              <img className="user-avatar" src={formatter.avatarUrl(this.props.data.user.avatar_thumbnail_url)} />
            </div>
            <div className="post-user">
              <h4 className="post-username">
                <a href={"/profiles/" + this.props.data.user.username}>@{this.props.data.user.username}</a>
              </h4>
              <span className="post-timestamp">{moment(this.props.data.post.created_at).fromNow()}</span>
            </div>
          </div>
          {attachment}
          <div className="post-body" dangerouslySetInnerHTML={formatter.post(this.props.data.post.body)} />
        </div>
        <div className="post-actions">
          <CommentCounter count={this.props.data.comment_count} className="right" onClick={this.handleToggleComments} />
          <Likes onClick={this.handleLike} count={this.props.data.like_count} liked={this.props.data.liked} />
          <a href="javascript:void(0)" onClick={this.handleFlag} className="button">
            <Icon name="flag" />
            Report
          </a>
          <div className="clearfix"></div>
        </div>
        <div className="post-comments">
          {comments}
        </div>
      </div>
    );
  }
});
