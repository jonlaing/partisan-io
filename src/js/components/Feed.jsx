import React from 'react/addons';

import Icon from 'react-fontawesome';

import FeedActionCreator from '../actions/FeedActionCreator.js';
import FeedStore from '../stores/FeedStore.js';

import FeedList from './FeedList.jsx';
import PostComposer from './PostComposer.jsx';
import ProfileEdit from './ProfileEdit.jsx';
import FlagForm from './FlagForm.jsx';
import UserSession from './UserSession.jsx';
import Nav from './Nav.jsx';
import MiniMatcher from './MiniMatcher.jsx';

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
      noFriends = (
        <div className="feed-nothing">
          <h3>You don't have any friends! <Icon name="frown-o"/></h3>
          <div>Well, at least not on Partisan. To find friends check out your matches, where we'll find people you'll probably vibe with</div>
          <a className="button" href="/matches">Find Matches</a>
        </div>
      );
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
