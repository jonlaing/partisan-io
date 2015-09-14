import React from 'react';
import LikeActionCreator from '../actions/LikeActionCreator';
import LikeStore from '../stores/LikeStore';

export default React.createClass({
  getInitialState() {
    return { likeCount: 0, liked: false };
  },

  handleLike() {
    this.setState({liked: !this.state.liked});
    LikeActionCreator.like(this.props.type, this.props.id);
  },

  componentDidMount() {
    LikeStore.addChangeListener(this._onChange);
    LikeActionCreator.getLikes(this.props.type, this.props.id);
  },

  componentWillUnmount() {
    LikeStore.removeChangeListener(this._onChange);
  },

  render() {
    return (
      <div>
        <button className={"like" + (this.state.liked ? " active" : "")} onClick={this.handleLike}>
          <i className="fi-like"></i>
          {this.state.likeCount}
        </button>
      </div>
    );
  },

  _onChange() {
    let state = LikeStore.getLikes(this.props.type, this.props.id);
    this.setState(state);
  }
});
