import React from 'react/addons';
import FeedActionCreator from '../actions/FeedActionCreator.js';
import FeedStore from '../stores/FeedStore.js';

// import ProfileEdit from './ProfileEdit.jsx';
import Card from './Card.jsx';
import Post from './Post.jsx';
import PostComposer from './PostComposer.jsx';
import FlagForm from './FlagForm.jsx';
import UserSession from './UserSession.jsx';
import Nav from './Nav.jsx';

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
      nothing = (<strong>Nothing here yet!</strong>);
    }

    return (
      <div className="feed">
        <header>
          <UserSession className="right" username={this.props.data.user.username} avatar={this.props.data.user.avatar_thumbnail_url} />
          <img src="images/logo.svg" className="logo" />
          <Nav currentPage="feed" />
        </header>

        <div className="container">
          <aside>
          </aside>
          <article>
            <PostComposer />
            <ReactCSSTransitionGroup transitionName="feed">
              {cards}
            </ReactCSSTransitionGroup>
            {nothing}
          </article>
          <aside>&nbsp;</aside>
        </div>

        <FlagForm show={this.state.modals.flag.show} id={this.state.modals.flag.id} type={this.state.modals.flag.type} ref="flag"/>
      </div>
    );
  },

  _onChange() {
    this.setState(FeedStore.getState());
  }
});
