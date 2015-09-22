import React from 'react/addons';

import CommentsActionCreator from '../actions/CommentsActionCreator';
import FeedStore from '../stores/FeedStore';
import Comment from './Comment.jsx';

var ReactCSSTransitionGroup = React.addons.CSSTransitionGroup;

// when we go to a page with a comment hash, scroll to it, but only once
var _scrolled = false;

export default React.createClass({
  getInitialState() {
    return {comments: []};
  },

  componentDidMount() {
    FeedStore.addChangeListener(this._onChange);
    CommentsActionCreator.getList(this.props.type, this.props.id);
  },

  componentWillUnmount() {
    FeedStore.removeChangeListener(this._onChange);
  },

  componentDidUpdate() {
    if(window.location.hash && _scrolled === false) {
      $('html, body').animate({
        scrollTop: $(window.location.hash).offset().top
      });
      _scrolled = true;
    }
  },

  render() {
    var comments;

    comments = this.state.comments.map(function(comment) {
      let hash = "comment-" + comment.comment.id;
      return (
        <li key={comment.comment.id}>
          <a id={hash} name={hash}></a>
          <Comment data={comment} />
        </li>
      );
    });

    return (
      <div className="commentlist">
        <ReactCSSTransitionGroup transitionName="commentlist-item" component="ul">
          {comments}
        </ReactCSSTransitionGroup>
      </div>
    );
  },

  _onChange() {
    this.setState({comments: FeedStore.listComments(this.props.type, this.props.id)});
  }
});
