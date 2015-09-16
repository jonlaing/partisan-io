import React from 'react';
import moment from 'moment';

import Likes from './Likes.jsx';
import CommentCounter from './CommentCounter.jsx';
import CommentComposer from './CommentComposer.jsx';
import CommentList from './CommentList.jsx';

export default React.createClass({
  getInitialState() {
    return {showComments: false};
  },

  handleToggleComments() {
    let show = !this.state.showComments;
    this.setState({showComments: show});
  },

  render() {
    var comments;

    if(this.state.showComments === true) {
      comments = (
        <div>
          <CommentList id={this.props.data.post.id} type="posts" show={this.state.showComments}/>
          <CommentComposer id={this.props.data.post.id} type="posts" />
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
          <div className="post-body">
            {this.props.data.post.body}
          </div>
        </div>
        <div className="post-actions">
          <CommentCounter id={this.props.data.post.id} type="posts" className="right" onClick={this.handleToggleComments} />
          <Likes id={this.props.data.post.id} type="posts" />
          <div className="clearfix"></div>
        </div>
        <div className="post-comments">
          {comments}
        </div>
      </div>
    );
  }
});
