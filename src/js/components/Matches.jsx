import React from 'react';

import moment from 'moment';

import MatchesActionCreator from '../actions/MatchesActionCreator';
import MatchesStore from '../stores/MatchesStore';

import UserSession from './UserSession.jsx';
import Nav from './Nav.jsx';
import Match from './Match.jsx';
import MatchSearch from './MatchSearch.jsx';

export default React.createClass({
  getInitialState() {
    return {matches: []};
  },

  handleMatchClick(username) {
    return () => {
      window.location.href = '/profiles/' + username;
    };
  },

  componentDidMount() {
    MatchesStore.addChangeListener(this._onChange);
    MatchesActionCreator.getMatches(25, "", 18, moment().diff(this.props.data.user.birthdate, 'years') + 10);
  },

  componentWillUnount() {
    MatchesStore.removeChangeListener(this._onChange);
  },

  render() {
    var nothing, matches;
    let self = this;

    matches = this.state.matches.map(function(match, i) {
      return (
        <li key={i} onClick={self.handleMatchClick(match.user.username)}>
          <Match user={match.user} match={match.match} />
        </li>
      );
    });

    if(this.state.matches.length < 1) {
      nothing = <span>You have no matches</span>;
    }

    return (
      <div className="matches">
        <header className="header">
          <UserSession className="right" username={this.props.data.user.username} avatar={this.props.data.user.avatar_thumbnail_url} />
          <img src="images/logo.svg" className="logo" />
          <Nav currentPage="matches" />
        </header>

        <div className="dashboard">
          <aside>
            <MatchSearch birthdate={this.props.data.user.birthdate} />
          </aside>
          <article className="matches-container">
            <ul className="matchlist">
              {matches}
            </ul>
            {nothing}
          </article>
        </div>
      </div>
    );
  },

  _onChange() {
    let state = MatchesStore.getAll();
    this.setState(state);
  },

  _matchClass() {
  }
});
