import React from 'react';

import FriendsActionCreator from '../actions/FriendsActionCreator';
import FriendsStore from '../stores/FriendsStore';

import UserSession from './UserSession.jsx';
import Nav from './Nav.jsx';
import Match from './Match.jsx';
import NoFriends from './NoFriends.jsx';
import MiniMatcher from './MiniMatcher.jsx';

export default React.createClass({
  getInitialState() {
    return { 
      friendships: [],
      filter : ""
    };
  },

  handleFriendClick(username) {
    return () => {
      window.location.href = '/profiles/' + username;
    };
  },

  handleSearchChange(e) {
    this.setState({filter: e.target.value});
  },

  componentDidMount() {
    FriendsStore.addChangeListener(this._onChange);
    FriendsActionCreator.getAll();
  },

  componentWillUnmount() {
    FriendsStore.removeChangeListener(this._onChange);
  },

  render() {
    var noFriends;

    var self = this;

    var friendlist = this.state.friendships.filter((f) => f.user.username.includes(this.state.filter)).map((friendship, i) => {
      return (
        <li key={i} onClick={self.handleFriendClick(friendship.user.username)}>
          <Match user={friendship.user} match={friendship.match} />
        </li>
      );
    });

    if(this.state.friendships.length === 0) {
      noFriends = <NoFriends />;
    }

    return (
      <div className="friends">
        <header className="header">
          <UserSession className="right" username={this.props.data.user.username} avatar={this.props.data.user.avatar_thumbnail_url} />
          <img src="images/logo.svg" className="logo" />
          <Nav currentPage="friends" />
        </header>

        <div className="dashboard">
          <article className="friends-container">
            <div className="friends-search search">
              <input type="text" placeholder="Search for friends…" onChange={this.handleSearchChange} />
            </div>
            <ul className="friendlist">
              {friendlist}
            </ul>
            {noFriends}
          </article>
          <aside>
            <MiniMatcher />
          </aside>
        </div>
      </div>
    );
  },

  _onChange() {
    let state = FriendsStore.getAll();
    this.setState(state);
  }
});
