import React from 'react';

export default React.createClass({
  getInitialState() {
    return {matches: []};
  },

  componentDidMount() {
  },

  render() {
    var matches = this.state.matches.map(function(match, i) {
      return (
        <li key={i}>
          <div class="row">
            <div class="large-6 columns">
              <a href={"profiles/" + match.user.id}>@{match.user.username}</a>
            </div>
            <div class="large-3 columns">
              {match.match}% Match
              {match.enemy}% Enemy
            </div>
          </div>
        </li>
      );
    });


    return (
      <div className="matches">
        <ul>
          {matches}
        </ul>
      </div>
    );
  }
});
