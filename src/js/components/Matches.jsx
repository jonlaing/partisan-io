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
    var nothing, matches;

    matches = this.state.matches.map(function(match, i) {
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

    if(this.state.matches.length < 1) {
      nothing = <span>You have no matches</span>;
    }

    return (
      <div className="matches">
        <ul>
          {matches}
        </ul>
        {nothing}
      </div>
    );
  },

  _onChange() {
    let state = MatchesStore.getAll();
    this.setState(state);
  }
});
