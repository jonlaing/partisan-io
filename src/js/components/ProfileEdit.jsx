import React from 'react';

import CheckboxGroup from 'react-checkbox-group';
import Icon from 'react-fontawesome';

import formatter from '../utils/formatter';

import AvatarEditor from './AvatarEditor.jsx';
import ProfileInfoEditor from './ProfileInfoEditor.jsx';

import ProfileActionCreator from '../actions/ProfileActionCreator';
import ProfileStore from '../stores/ProfileStore';


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

  handleAvatarFinish(avatar) {
    this.setState({avatarUrl: avatar.thumb});
  },

  handleLocationFinish(e) {
      ProfileActionCreator.updateLocation(e.target.value);
  },

  handleGenderFinish(e) {
      ProfileActionCreator.updateGender(e.target.value);
  },

  handleLookingForChange() {
    var bitMap = 0;
    let values = this.refs.lookingFor.getCheckedValues();

    values.forEach(function(v) {
      bitMap += parseInt(v);
    });

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
    var summary;

    if(this.state.editSummary === false) {
      summary = (
        <div>
          <h3>
            Summary
            <Icon name="pencil" onClick={this.handleSummaryClick} className="profile-summary-edit"/>
          </h3>
          <div dangerouslySetInnerHTML={ formatter.userSummary(this.state.profile.summary) } />
        </div>
      );
    } else {
      summary = this._editSummaryTemplate();
    }

    return (
      <div className="profile-edit">
        <div className="profile-edit-container">
          <div className="profile-avatar-container">
            <AvatarEditor onSuccess={this.handleAvatarFinish} avatarUrl={this.state.avatarUrl} />
          </div>
          <h2 className="profile-username">
            @{this.props.data.user.username}
          </h2>
          <ProfileInfoEditor
            location={this.state.user.location}
            gender={this.state.user.gender}
            postalCode={this.state.user.postal_code}
            onLocationFinish={this.handleLocationFinish}
            onGenderFinish={this.handleGenderFinish} />
          <div className="profile-edit-lookingfor">
            <h3>Looking For</h3>
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
          <div className="profile-edit-summary">
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

  _editSummaryTemplate() {
    return (
      <div>
        <h3>
          Summary
        </h3>
        <textarea defaultValue={this.state.profile.summary} placeholder="Tell us a little about yourself!" ref="summary" />
        <button onClick={this.handleSummarySubmit} >Done</button>
      </div>
    );
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
