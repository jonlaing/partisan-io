import React from 'react/addons';

import Card from './Card.jsx';
import Post from './Post.jsx';

var ReactCSSTransitionGroup = React.addons.CSSTransitionGroup;

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
        <ReactCSSTransitionGroup transitionName="feed">
          {cards}
        </ReactCSSTransitionGroup>
      </div>
    );
  }
});
