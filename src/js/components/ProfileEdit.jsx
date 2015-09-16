import React from 'react';

import AvatarUpload from './AvatarUpload.jsx';

import ProfileActionCreator from '../actions/ProfileActionCreator';
import ProfileStore from '../stores/ProfileStore';

let _ENTER = 13;

export default React.createClass({
  getInitialState() {
    return {
      showAvatarUpload: false,
      editLocation: false,
      editGender: false,
      avatarUrl: this.props.data.user.avatar_thumbnail_url,
      user: this.props.data.user,
      profile: this.props.data.profile
    };
  },

  handleAvatarClick() {
    this.setState({showAvatarUpload: true});
  },

  handleAvatarFinish(avatar) {
    this.setState({showAvatarUpload: false, avatarUrl: avatar.thumb});
  },

  handleAvatarCancel() {
    this.setState({showAvatarUpload: false});
  },

  handleLocationClick() {
    this.setState({editLocation: true});
  },

  handleLocationKeyDown(e) {
    if(e.keyCode === _ENTER && e.target.value !== "") {
      this.setState({editLocation: false});
      ProfileActionCreator.updateLocation(e.target.value);
    }
  },

  handleGenderClick() {
    this.setState({editGender: true});
  },

  handleGenderKeyDown(e) {
    if(e.keyCode === _ENTER && e.target.value !== "") {
      this.setState({editGender: false});
      ProfileActionCreator.updateGender(e.target.value);
    }
  },

  componentDidMount() {
    ProfileStore.addChangeListener(this._onChange);
  },

  componentWillUnmount() {
    ProfileStore.removeChangeListener(this._onChange);
  },

  render() {
    var avatar, location, gender;

    if(this.state.showAvatarUpload === false) {
      avatar = <img src={this.state.avatarUrl} onClick={this.handleAvatarClick} />;
    } else {
      avatar = this._uploadTemplate();
    }

    if(this.state.editLocation === false) {
      var cityState = this._cityState(this.state.user.location);

      location = (
        <div className="large-10 columns">
          <a href="javascript:void(0)" onClick={this.handleLocationClick}>{cityState}</a>
        </div>
      );
    } else {
      location = this._editLocationTemplate();
    }

    if(this.state.editGender === false) {
      let g = this.state.user.gender || "None";
      gender = (
        <div className="large-10 columns">
          <a href="javascript:void(0)" onClick={this.handleGenderClick}>{g}</a>
        </div>
      );
    } else {
      gender = this._editGenderTemplate();
    }

    return (
      <div className="profile-edit">
        <div className="row">
          <div className="large-3 columns">
            <div className="profile-avatar">
              <div className="user-avatar">
                {avatar}
              </div>
            </div>
          </div>
          <div className="large-9 columns">
            <h1 className="profile-username">
              @{this.props.data.user.username}
            </h1>
          </div>
        </div>
        <div className="row">
          <div className="large-6 columns">
            <div className="row">
              <div className="large-2 columns">
                Location
              </div>
              {location}
            </div>
          </div>
          <div className="large-6 columns">
            <div className="row">
              <div className="large-2 columns">
                Gender
              </div>
              {gender}
            </div>
          </div>
        </div>
      </div>
    );
  },

  _onChange() {
    let state = ProfileStore.getProfile();
    this.setState(state);
  },

  _uploadTemplate() {
    return (
      <div>
        <AvatarUpload onSuccess={this.handleAvatarFinish} />
        <br/>
        <a href="javascript:void(0)" onClick={this.handleAvatarCancel}>Cancel</a>
      </div>
    );
  },
  _editLocationTemplate() {
    return (
      <div className="large-10 columns">
        <input type="text" defaultValue={this.state.user.postal_code} onKeyDown={this.handleLocationKeyDown} />
      </div>
    );
  },
  _editGenderTemplate() {
    let g = this.state.user.gender || "";
    return (
      <div className="large-10 columns">
        <input type="text" placeholder="Type in your gender" defaultValue={g} onKeyDown={this.handleGenderKeyDown} />
      </div>
    );
  },
  _cityState(location) {
    return location.replace(/\s\d+.*$/, '');
  }
});
