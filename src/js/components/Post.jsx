import React from 'react';
import moment from 'moment';
import marked from 'marked';

import LikeActionCreator from '../actions/LikeActionCreator';
import FlagActionCreator from '../actions/FlagActionCreator';

import Likes from './Likes.jsx';
import CommentCounter from './CommentCounter.jsx';
import CommentComposer from './CommentComposer.jsx';
import CommentList from './CommentList.jsx';

marked.setOptions({
  sanitize: true,
  tables: false
});

var markedRenderer = new marked.Renderer();

markedRenderer.heading = (text) => {
  return '<p><strong>' + text + '</strong></p>';
};

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

  render() {
    var comments, attachment;

    var bodyHTML = function(body) {
      return { __html: _hashtagify(marked(body, {renderer: markedRenderer} )) };
    };

    if(this.state.showComments === true) {
      comments = (
        <div>
          <CommentList id={this.props.data.post.id} show={this.state.showComments}/>
          <CommentComposer id={this.props.data.post.id} />
        </div>
      );
    }

    if(this.props.data.image_attachment.id > 0) {
      attachment = (
        <div className="post-attachment">
          <img src={this.props.data.image_attachment.image_url} width="100%" />
        </div>
      );
    }

    return (
      <div className="post">
        <div className="card-body">
          <div className="post-header">
            <div className="post-avatar">
              <img className="user-avatar" src={this.props.data.user.avatar_thumbnail_url} />
            </div>
            <div className="post-user">
              <h4 className="post-username">
                <a href={"/profiles/" + this.props.data.user.id}>@{this.props.data.user.username}</a>
              </h4>
              <span className="post-timestamp">{moment(this.props.data.post.created_at).fromNow()}</span>
            </div>
          </div>
          {attachment}
          <div className="post-body" dangerouslySetInnerHTML={bodyHTML(this.props.data.post.body)} />
        </div>
        <div className="post-actions">
          <CommentCounter count={this.props.data.comment_count} className="right" onClick={this.handleToggleComments} />
          <Likes onClick={this.handleLike} count={this.props.data.like_count} liked={this.props.data.liked} />
          <a href="javascript:void(0)" onClick={this.handleFlag}><i className="fi-flag"></i></a>
          <div className="clearfix"></div>
        </div>
        <div className="post-comments">
          {comments}
        </div>
      </div>
    );
  }
});

function _hashtagify(content) {
  let tags = content.match(/#[a-zA-Z]+/g);
  var newContent = content;

  if(tags !== null) {
    tags.forEach((tag) => {
      let encoded = encodeURIComponent(tag);
      newContent = newContent.replace(tag, '<a href="/hashtags?q=' + encoded + '">' + tag + '</a>');
    });
  }

  return newContent;
}
