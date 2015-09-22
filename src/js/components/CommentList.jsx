import React from 'react/addons';

import CommentsActionCreator from '../actions/CommentsActionCreator';
import FeedStore from '../stores/FeedStore';
import Comment from './Comment.jsx';

var ReactCSSTransitionGroup = React.addons.CSSTransitionGroup;

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

  render() {
    var comments;

    comments = this.state.comments.map(function(comment, i) {
      return (<li key={i}><Comment data={comment} /></li>);
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
