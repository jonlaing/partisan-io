/*global $ */
import React from 'react';
import PostActionCreator from '../actions/PostActionCreator';

export default React.createClass({
  handleCreate() {
    let body = $(React.findDOMNode(this.refs.body));
    if(body.val().length > 0) {
      PostActionCreator.createPost(body.val());
    }
    body.val('');
    body.removeClass('focus');
  },

  handleFocus() {
    let body = $(React.findDOMNode(this.refs.body));
    body.addClass('focus');
  },

  handleBlur() {
    let body = $(React.findDOMNode(this.refs.body));
    body.removeClass('focus');
  },

  render() {
    return (
      <div className="post-composer">
        <div className="post-composer-field">
          <textarea rows="1" placeholder="Write a new post" onFocus={this.handleFocus} ref="body"></textarea>
        </div>
        <div className="post-composer-actions clearfix">
          <button className="button right" onClick={this.handleCreate}>Post</button>
        </div>
      </div>
    );
  }
});
