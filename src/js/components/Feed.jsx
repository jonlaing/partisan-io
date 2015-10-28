import React from 'react/addons';

import FeedActionCreator from '../actions/FeedActionCreator.js';
import FeedStore from '../stores/FeedStore.js';

import FeedList from './FeedList.jsx';
import PostComposer from './PostComposer.jsx';
import ProfileEdit from './ProfileEdit.jsx';
import FlagForm from './FlagForm.jsx';
import UserSession from './UserSession.jsx';
import Nav from './Nav.jsx';
import MiniMatcher from './MiniMatcher.jsx';
import NoFriends from './NoFriends.jsx';

export default React.createClass({
  getInitialState() {
    return {
      feed: [],
      noFriends: false,
      modals: {
        flag: { show: false, flagID: 0 }
      }
    };
  },

  componentDidMount() {
    FeedStore.addChangeListener(this._onChange);
    FeedActionCreator.getFeed();
  },

  componentWillUnmount() {
    FeedStore.removeChangeListener(this._onChange);
  },

  render() {
    var noFriends;


    if(this.state.noFriends === true) {
      noFriends = <NoFriends />;
    }

    return (
      <div className="feed">
        <header>
          <UserSession className="right" username={this.props.data.user.username} avatar={this.props.data.user.avatar_thumbnail_url} />
          <img src="/images/logo.svg" className="logo" />
          <Nav currentPage="feed" />
        </header>

        <div className="dashboard dashboard-3col">
          <aside>
            <ProfileEdit data={this.props.data} />
          </aside>
          <article>
            <PostComposer />
            <FeedList feed={this.state.feed} noFriends={this.state.noFriends} />
            {noFriends}
          </article>
          <aside>
            <MiniMatcher />
          </aside>
        </div>

        <FlagForm show={this.state.modals.flag.show} id={this.state.modals.flag.id} type={this.state.modals.flag.type} ref="flag"/>
      </div>
    );
  },

  _onChange() {
    this.setState(FeedStore.getState());
  }
});
