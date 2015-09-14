import React from 'react/addons';

import CommentsActionCreator from '../actions/CommentsActionCreator';
import CommentsStore from '../stores/CommentsStore';
import Comment from './Comment.jsx';

var ReactCSSTransitionGroup = React.addons.CSSTransitionGroup;

export default React.createClass({
  getInitialState() {
    return {comments: []};
  },

  componentDidMount() {
    CommentsStore.addChangeListener(this._onChange);
    CommentsActionCreator.getList(this.props.type, this.props.id);
  },

  componentWillUnmount() {
    CommentsStore.removeChangeListener(this._onChange);
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
    this.setState({comments: CommentsStore.list(this.props.type, this.props.id)});
  }
});
