import React from 'react';
import formatter from '../utils/formatter';

import MatchesActionCreator from '../actions/MatchesActionCreator';
import MatchesStore from '../stores/MatchesStore';

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

  componentWillUnmount() {
    MatchesStore.removeChangeListener(this._onChange);
  },

  render() {
    var matches = this.state.matches.map((match) => {
      var ageContent, gContent;

      var age = formatter.age(match.user.birthdate, false);
      var g = match.user.gender;

      if(age.length > 0) {
        ageContent = <span>{age}&nbsp;&bull;&nbsp;</span>;
      }

      if(g.length > 0) {
        gContent = <span>{g}&nbsp;&bull;&nbsp;</span>;
      }

      return (
        <li className="matchlist-small-item" key={match.user.id}>
          <div className="matchlist-small-avatar" onClick={this.handleAvatarClick(match.user.username)}>
            <img className="user-avatar" src={formatter.avatarUrl(match.user.avatar_thumbnail_url)} />
          </div>
          <div className="matchlist-small-user">
            <div className="matchlist-small-username">
              <a href={"/profiles/" + match.user.username}>@{match.user.username}</a>
            </div>
            <div className="matchlist-small-info">
              {ageContent}
              {gContent}
              {formatter.cityState(match.user.location)}
            </div>
          </div>
          <div className="matchlist-small-match">{match.match}%</div>
        </li>
      );
    });

    return (
      <div className="matchlist-small">
        <h3>People with similar views</h3>
        <ul>
          {matches}
        </ul>
        <div className="more">
          <a className="button" href="/matches">Find More Matches</a>
        </div>
      </div>
    );
  },

  _onChange() {
    let state = MatchesStore.getAll();
    this.setState(state);
  }
});
