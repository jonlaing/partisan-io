import React from 'react';
import ReactDOM from 'react-dom';

import Icon from 'react-fontawesome';

import formatter from '../utils/formatter';

import AvatarEditor from './AvatarEditor.jsx';
import ProfileInfoEditor from './ProfileInfoEditor.jsx';
import LookingForEdit from './LookingForEdit.jsx';

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

  handleBirthdateFinish(date) {
      ProfileActionCreator.updateBirthdate(date);
  },

  handleLookingForChange(val) {
    ProfileActionCreator.updateLookingFor(val);
  },

  handleSummaryClick() {
    this.setState({editSummary: true});
  },

  handleSummarySubmit() {
    let summary = $(ReactDOM.findDOMNode(this.refs.summary)).val();

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
      var s;
      if(this.state.profile.summary.length < 1) {
        s = (<div onClick={this.handleSummaryClick}><em>You haven&apos;t filled out your summary yet</em></div>);
      } else {
        s = (<div className="profile-summary-text" dangerouslySetInnerHTML={ formatter.userSummary(this.state.profile.summary) } />);
      }

      summary = (
        <div>
          <h3>
            Summary
            <Icon name="pencil" onClick={this.handleSummaryClick} className="profile-summary-edit"/>
          </h3>
          {s}
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
          <div className="profile-lookingfor">
            <h3>Looking For:</h3>
            <LookingForEdit lookingFor={this.state.profile.looking_for} onChange={this.handleLookingForChange}/>
          </div>
          <ProfileInfoEditor
            location={this.state.user.location}
            gender={this.state.user.gender}
            birthdate={this.state.user.birthdate}
            postalCode={this.state.user.postal_code}
            onLocationFinish={this.handleLocationFinish}
            onGenderFinish={this.handleGenderFinish}
            onBirthdateFinish={this.handleBirthdateFinish} />
          <div className="profile-summary">
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
        <div className="profile-edit-summary">
          <textarea defaultValue={this.state.profile.summary} placeholder="Tell us a little about yourself!" ref="summary" autoFocus={true} />
          <button className="button-small" onClick={this.handleSummarySubmit} >Done</button>
        </div>
      </div>
    );
  }

});
