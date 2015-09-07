import React from 'react/addons';
import FeedActionCreator from '../actions/FeedActionCreator.js';
import FeedStore from '../stores/FeedStore.js';

import Card from './Card.jsx';
import Post from './Post.jsx';
import PostComposer from './PostComposer.jsx';
import UserSession from './UserSession.jsx';

var ReactCSSTransitionGroup = React.addons.CSSTransitionGroup;

function getStateFromStore() {
  let state = FeedStore.getAll();
  return state;
}

export default React.createClass({
  getInitialState() {
    return { feed: [] };
  },

  componentDidMount() {
    FeedStore.addChangeListener(this._onChange);
    FeedActionCreator.getFeed();
  },

  componentWillUnmount() {
    FeedStore.removeChangeListener(this._onChange);
  },

  render() {
    var cards;

    cards = this.state.feed.map(function(item, i) {
      return (
        <Card key={i}>
          <Post data={item.record} />
        </Card>
      );
    });

    return (
      <div className="feed">
        <PostComposer />
        <ReactCSSTransitionGroup transitionName="feed">
          {cards}
        </ReactCSSTransitionGroup>
      </div>
    );
  },

  _onChange() {
    this.setState({feed: getStateFromStore()});
  }
});
