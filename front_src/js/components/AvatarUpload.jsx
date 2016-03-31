import React from 'react';

import Icon from 'react-fontawesome';
import Dropzone from 'react-dropzone';

import AvatarActionCreator from '../actions/AvatarActionCreator';
import AvatarStore from '../stores/AvatarStore';

export default React.createClass({
  getInitialState() {
    return {files: [], avatar: {}};
  },

  handleDrop(files) {
    this.setState({files: files});
    AvatarActionCreator.uploadAvatar(files);
  },

  componentDidMount() {
    AvatarStore.addChangeListener(this._onChange);
  },

  componentWillUnmount() {
    AvatarStore.removeChangeListener(this._onChange);
  },

  componentDidUpdate() {
    if(this.state.avatar.full !== undefined) {
      this.props.onSuccess(this.state.avatar);
    }
  },

  render() {
    var content;

    if(this.state.files.length < 1) {
      content = this._dropTemplate();
    } else {
      content = this._previewTemplate();
    }

    return (
      <div className="avatar-uploader">
        {content}
      </div>
    );
  },

  _dropTemplate() {
    return (
      <Dropzone multiple={false} onDrop={this.handleDrop} className="dropzone" activeClassName="dropzone-active">
        <div>
          <Icon name="download" />
          <span className="help-text">Drag image here, or click to browse</span>
        </div>
      </Dropzone>
    );
  },

  _previewTemplate() {
    return (
      <div>
        <i className="fa fa-circle-o-notch fa-spin"></i>
        &nbsp;Uploading&hellip;
      </div>
    );
  },

  _onChange() {
    let state = AvatarStore.getAvatar();
    this.setState({avatar: state});
  }
});
