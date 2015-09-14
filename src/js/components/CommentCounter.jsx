import React from 'react';

import CommentsActionCreator from '../actions/CommentsActionCreator';
import CommentsStore from '../stores/CommentsStore';

export default React.createClass({
  getInitialState() {
    return {commentCount: 0};
  },

  componentDidMount() {
    CommentsStore.addChangeListener(this._onChange);
    CommentsActionCreator.getCount(this.props.type, this.props.id);
  },

  componentWillUnmount() {
    CommentsStore.removeChangeListener(this._onChange);
  },

  render() {
    return (
      <div className={this.props.className}>
        <button className="comment" onClick={this.props.onClick} >
          <i className="fi-comment"></i>
          {this.state.commentCount}
        </button>
      </div>
    );
  },

  _onChange() {
    let state = CommentsStore.count(this.props.type, this.props.id);
    this.setState({commentCount: state});
  }
});
