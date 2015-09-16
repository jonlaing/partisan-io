import React from 'react';

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
      <div>
        {content}
      </div>
    );
  },

  _dropTemplate() {
    return (
      <Dropzone multiple={false} onDrop={this.handleDrop}>
        <div className="text-center">
          Drop files here <br/>
          Or click here to browse
        </div>
      </Dropzone>
    );
  },

  _previewTemplate() {
    return (
      <div>
        Uploading ...
      </div>
    );
  },

  _onChange() {
    let state = AvatarStore.getAvatar();
    this.setState({avatar: state});
  }
});
