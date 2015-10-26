import React from 'react';

import formatter from '../utils/formatter';

import MatchesActionCreator from '../actions/MatchesActionCreator';
import MatchesStore from '../stores/MatchesStore';

import UserSession from './UserSession.jsx';
import Nav from './Nav.jsx';

export default React.createClass({
  getInitialState() {
    return {matches: []};
  },

  handleAvatarClick(username) {
    return () => {
      window.location.href = '/profiles/' + username;
    };
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
    let self = this;

    matches = this.state.matches.map(function(match, i) {
      return (
        <li key={i} onClick={self.handleAvatarClick(match.user.username)}>
          <div>
            <div className="matchlist-avatar">
              <img src={formatter.avatarUrl(match.user.avatar_thumbnail_url)} className="user-avatar" />
            </div>
            <div>
              <div className="matchlist-user">
                <a href={"profiles/" + match.user.username}>@{match.user.username}</a>
              </div>
              <div className="matchlist-info">{formatter.age(match.user.birthdate)}&nbsp;-&nbsp;{match.user.gender}</div>
              <div className="matchlist-location">{formatter.cityState(match.user.location)}</div>
              <div className="matchlist-match">{formatter.match(match.match)}</div>
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
        <header>
          <UserSession className="right" username={this.props.data.user.username} avatar={this.props.data.user.avatar_thumbnail_url} />
          <img src="images/logo.svg" className="logo" />
          <Nav currentPage="matches" />
        </header>

        <div className="dashboard">
          <aside>Blah</aside>
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
