import React from 'react';
import CheckboxGroup from 'react-checkbox-group';

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
      editSummary: false,
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

  handleLookingForChange() {
    var bitMap = 0;
    let values = this.refs.lookingFor.getCheckedValues();

    values.forEach(function(v) {
      bitMap += parseInt(v);
    });

    console.log(bitMap);

    ProfileActionCreator.updateLookingFor(bitMap);
  },

  handleSummaryClick() {
    this.setState({editSummary: true});
  },

  handleSummarySubmit() {
    let summary = $(React.findDOMNode(this.refs.summary)).val();

    this.setState({editSummary: false});
    ProfileActionCreator.updateSummary(summary);
  },

  componentDidMount() {
    ProfileStore.addChangeListener(this._onChange);
  },

  componentWillUnmount() {
    ProfileStore.removeChangeListener(this._onChange);
  },

  render() {
    var avatar, location, gender, summary;

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

    if(this.state.editSummary === false) {
      summary = (
        <div>
          <p>{this.state.profile.summary}</p>
          <a href="javascript:void(0)" onClick={this.handleSummaryClick}>Edit</a>
        </div>
      );
    } else {
      summary = this._editSummaryTemplate();
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
          <div className="large-4 columns">
            <CheckboxGroup
              name="looking_for"
              value={this._parseLookingFor(this.state.profile.looking_for)}
              ref="lookingFor"
              onChange={this.handleLookingForChange} >
                <label>
                  <input type="checkbox" value={1 << 0} /> Friends
                </label>
                <label>
                  <input type="checkbox" value={1 << 1} /> Love
                </label>
                <label>
                  <input type="checkbox" value={1 << 2} /> Enemy
                </label>
            </CheckboxGroup>
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
        <div className="row">
          <div className="large-12">
            <h4>Summary</h4>
            {summary}
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
  _editSummaryTemplate() {
    return (
      <div>
        <textarea defaultValue={this.state.profile.summary} placeholder="Tell us a little about yourself!" ref="summary" />
        <button onClick={this.handleSummarySubmit} >Save</button>
      </div>
    );
  },

  _cityState(location) {
    return location.replace(/\s\d+.*$/, '');
  },
  _parseLookingFor(n) {
    var vals = [];
    for(var i = 0; i <= 3; i++) {
      if((n & 1 << i) !== 0) {
        vals.push((1 << i).toString());
      }
    }
    return vals;
  }
});
