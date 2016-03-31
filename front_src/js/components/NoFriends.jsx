import React from 'react';

import Icon from 'react-fontawesome';

export default React.createClass({
  getInitialState() {
    return {};
  },

  componentDidMount() {
  },

  render() {
    return (
      <div className="friends-none">
        <h3>You don't have any friends! <Icon name="frown-o"/></h3>
        <div>Well, at least not on Partisan. To find friends check out your matches, where we'll find people you'll probably vibe with</div>
        <a className="button" href="/matches">Find Matches</a>
      </div>
    );
  }
});
