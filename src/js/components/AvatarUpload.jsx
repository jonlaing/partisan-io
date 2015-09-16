import React from 'react';

import Dropzone from 'react-dropzone';
import AvatarActionCreator from '../actions/AvatarActionCreator';

export default React.createClass({
  getInitialState() {
    return {files: []};
  },

  handleDrop(files) {
    this.setState({files: files});
    AvatarActionCreator.uploadAvatar(files);
  },

  componentDidMount() {
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
        Drop your file here
      </Dropzone>
    );
  },

  _previewTemplate() {
    return (
      <div>
        <div>
          <img src={this.state.files[0].preview} width="100"/>
        </div>
        <div>
          Uploading ...
        </div>
      </div>
    );
  }
});
