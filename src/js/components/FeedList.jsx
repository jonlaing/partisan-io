import React from 'react';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';

import Card from './Card.jsx';
import Post from './Post.jsx';

export default React.createClass({
  render() {
    var cards;

    cards = this.props.feed.map(function(item, i) {
      if(item.record_type === "post") {
        return (
          <Card key={i}>
            <Post data={item.record} />
          </Card>
        );
      }
    });

    return (
      <div>
        <ReactCSSTransitionGroup transitionName="feed" transitionEnterTimeout={1000} transitionLeaveTimeout={1000}>
          {cards}
        </ReactCSSTransitionGroup>
      </div>
    );
  }
});
