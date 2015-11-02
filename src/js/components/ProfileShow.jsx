import React from 'react';

import formatter from '../utils/formatter';

import FeedActionCreator from '../actions/FeedActionCreator.js';
import FeedStore from '../stores/FeedStore.js';

import FeedList from './FeedList.jsx';
import Friender from './Friender.jsx';
import UserSession from './UserSession.jsx';
import Nav from './Nav.jsx';
import FlagForm from './FlagForm.jsx';

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
    FeedActionCreator.getByUser(this.props.user.id);
  },

  componentWillUnmount() {
    FeedStore.removeChangeListener(this._onChange);
  },

  render() {
    var gender;

    if(this.props.user.gender !== undefined) {
      gender = <div className="profile-info-gender">{this.props.user.gender}</div>;
    }

    return (
      <div className="profile">
        <header className="header">
          <UserSession className="right" username={this.props.currentUser.username} avatar={formatter.avatarUrl(this.props.currentUser.avatar_thumbnail_url)} />
          <img src="/images/logo.svg" className="logo" />
          <Nav />
        </header>

        <div className="dashboard dashboard-3col">
          <aside>
            <div className="profile-show-container">
              <div className="profile-avatar-container">
                <div className="profile-avatar">
                  <img src={formatter.avatarUrl(this.props.user.avatar_thumbnail_url)} className="user-avatar" />
                </div>
              </div>
              <h2 className="profile-username">@{this.props.user.username}</h2>
              <div className="profile-lookingfor">
                <h3>Looking For</h3>
                <div className="lookingfor">
                  <div>
                    <label className={(this._active(1 << 0) ? " active" : "")}>
                      <i className="fi-torsos-all-female"></i>
                      Friends
                    </label>
                    <label className={(this._active(1 << 1) ? " active" : "")}>
                      <i className="fi-heart"></i>
                      Love
                    </label>
                    <label className={(this._active(1 << 2) ? " active" : "")}>
                      <i className="fi-skull"></i>
                      Enemies
                    </label>
                  </div>
                </div>
              </div>
              <div className="profile-match-container">
                <div className="profile-match">
                  {this.props.match}% Match
                </div>
              </div>
              <div className="profile-info">
                <div className="profile-info-age">{formatter.age(this.props.user.birthdate)}</div>
                <div className="profile-info-location">{formatter.cityState(this.props.user.location)}</div>
                {gender}
              </div>
              <div className="profile-summary">
                <h3>Summary</h3>
                <div className="profile-summary-text" dangerouslySetInnerHTML={formatter.userSummary(this.props.profile.summary)} />
              </div>
            </div>
          </aside>
          <article>
            <div>
              <Friender id={this.props.user.id} username={this.props.user.username} />
            </div>
            <FeedList feed={this.state.feed} />
          </article>
          <aside>Something else goes here</aside>
        </div>

        <FlagForm show={this.state.modals.flag.show} id={this.state.modals.flag.id} type={this.state.modals.flag.type} ref="flag"/>
      </div>
    );
  },

  _onChange() {
    this.setState(FeedStore.getState());
  },

  _parseLookingFor(n) {
    var vals = [];
    for(var i = 0; i <= 3; i++) {
      if((n & 1 << i) !== 0) {
        vals.push((1 << i).toString());
      }
    }
    return vals;
  },

  _active(n) {
    let vals = this._parseLookingFor(this.props.profile.looking_for);
    for(let i in vals) {
      if(parseInt(vals[i]) === n) {
        return true;
      }
    }
    return false;
  }
});
