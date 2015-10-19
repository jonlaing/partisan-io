import React from 'react/addons';

import Icon from 'react-fontawesome';

import FeedActionCreator from '../actions/FeedActionCreator.js';
import FeedStore from '../stores/FeedStore.js';

import Card from './Card.jsx';
import Post from './Post.jsx';
import PostComposer from './PostComposer.jsx';
import ProfileEdit from './ProfileEdit.jsx';
import FlagForm from './FlagForm.jsx';
import UserSession from './UserSession.jsx';
import Nav from './Nav.jsx';
import MiniMatcher from './MiniMatcher.jsx';

var ReactCSSTransitionGroup = React.addons.CSSTransitionGroup;

export default React.createClass({
  getInitialState() {
    return {
      feed: [],
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
    var cards, nothing;

    cards = this.state.feed.map(function(item, i) {
      if(item.record_type === "post") {
        return (
          <Card key={i}>
            <Post data={item.record} />
          </Card>
        );
      }
    });

    if(this.state.feed.length === 0) {
      nothing = (
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
          <img src="images/logo.svg" className="logo" />
          <Nav currentPage="feed" />
        </header>

        <div className="container dashboard">
          <aside>
            <ProfileEdit data={this.props.data} />
          </aside>
          <article>
            <PostComposer />
            <ReactCSSTransitionGroup transitionName="feed">
              {cards}
            </ReactCSSTransitionGroup>
            {nothing}
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
