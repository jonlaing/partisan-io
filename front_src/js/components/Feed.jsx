import React from 'react';

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
import Lightbox from './Lightbox.jsx';

export default React.createClass({
  getInitialState() {
    return {
      feed: [],
      noFriends: false,
      scrollLoading: false,
      page: 1,
      modals: {
        flag: { show: false, flagID: 0 }
      }
    };
  },

  handleScroll() {
    var docHeight = $(document).height();
    var inHeight = window.innerHeight;
    var scrollT = $(window).scrollTop();
    var totalScrolled = scrollT + inHeight;
    if(totalScrolled + 100 > docHeight) {  //user reached at bottom
      if(this.state.scrollLoading === false) {  //to avoid multiple request
          this.setState({ scrollLoading: true, page: this.state.page + 1 });
          FeedActionCreator.getPage(this.state.page);
      }
    }
  },

  componentDidMount() {
    window.addEventListener("scroll", this.handleScroll);
    FeedStore.addChangeListener(this._onChange);
    FeedActionCreator.getFeed();
    FeedActionCreator.feedSocket();
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
        <header className="header">
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
        <Lightbox />
      </div>
    );
  },

  _onChange() {
    this.setState(FeedStore.getState());
  }
});
