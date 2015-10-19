import React from 'react';

import formatter from '../utils/formatter';

import AvatarUpload from './AvatarUpload.jsx';

export default React.createClass({
  getInitialState() {
    return {
      showAvatarUpload: false
    };
  },

  handleAvatarClick() {
    this.setState({showAvatarUpload: true});
  },

  handleAvatarCancel() {
    this.setState({showAvatarUpload: false});
  },

  handleSuccess(avatar) {
    this.setState({showAvatarUpload: false});
    this.props.onSuccess(avatar);
  },

  render() {
    if(this.state.showAvatarUpload === false) {
      return (
        <div className="profile-avatar" onClick={this.handleAvatarClick}>
          <img className="user-avatar" src={formatter.avatarUrl(this.props.avatarUrl)} />
        </div>
      );
    } else {
      return (
        <div>
          <AvatarUpload onSuccess={this.handleSuccess} />
          <br/>
          <a href="javascript:void(0)" onClick={this.handleAvatarCancel}>Cancel</a>
        </div>
      );
    }
  }
});
