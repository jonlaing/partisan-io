import React from 'react/addons';
import FeedActionCreator from '../actions/FeedActionCreator.js';
import FeedStore from '../stores/FeedStore.js';

import Card from './Card.jsx';
import Post from './Post.jsx';
import PostComposer from './PostComposer.jsx';

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
    var cards, nothing;

    cards = this.state.feed.map(function(item, i) {
      if(item.record_type === "post") {
        return (
          <Card key={i}>
            <Post data={item.record} />
          </Card>
        );
      }
    });

    if(this.state.feed.length === 0) {
      nothing = (<strong>Nothing here yet!</strong>);
    }

    return (
      <div className="feed">
        <PostComposer />
        <ReactCSSTransitionGroup transitionName="feed">
          {cards}
        </ReactCSSTransitionGroup>
        {nothing}
      </div>
    );
  },

  _onChange() {
    this.setState({feed: getStateFromStore()});
  }
});
