import React from 'react';

import MatchesActionCreator from '../actions/MatchesActionCreator';
import MatchesStore from '../stores/MatchesStore';

export default React.createClass({
  getInitialState() {
    return {matches: []};
  },

  componentDidMount() {
    MatchesStore.addChangeListener(this._onChange);
    MatchesActionCreator.getMatches();
  },

  componentWillUnount() {
    MatchesStore.removeChangeListener(this._onChange);
  },

  render() {
    var matches = this.state.matches.map(function(match, i) {
      return (
        <li key={i}>
          <div className="row">
            <div className="large-5 columns">
              <div>
                <a href={"profiles/" + match.user.id}>@{match.user.username}</a>
              </div>
              {match.user.location}
            </div>
            <div className="large-4 columns">
              {match.match}% Match
            </div>
            <div className="large-3 columns">
              <a className="button" href={"profiles/" + match.user.id}>View Profile</a>
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
  },

  _onChange() {
    this.setState(MatchesStore.getAll());
  }
});
