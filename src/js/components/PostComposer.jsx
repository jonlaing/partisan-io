/*global $ */
import React from 'react';
import Dropzone from 'react-dropzone';

import PostActionCreator from '../actions/PostActionCreator';

export default React.createClass({
  getInitialState() {
    return { showImageUploader: false, attachments: [] };
  },

  handleCreate() {
    let body = $(React.findDOMNode(this.refs.body));

    if(body.val().length > 0 || this.state.attachments.length > 0) {
      PostActionCreator.createPost(body.val(), this.state.attachments);
      body.val('');
      body.removeClass('focus');

      this.setState(this.getInitialState());
    }
  },

  handleFocus() {
    let body = $(React.findDOMNode(this.refs.body));
    body.addClass('focus');
  },

  handleBlur() {
    let body = $(React.findDOMNode(this.refs.body));
    body.removeClass('focus');
  },

  handleImageClick() {
    this.setState({showImageUploader: true});
  },

  handleDrop(files) {
    this.setState({attachments: files, showImageUploader: false});
  },

  handleImageCancel() {
    this.setState({attachments: [], showImageUploader: false});
  },

  render() {
    var imageUploader;

    if(this.state.showImageUploader === true) {
      imageUploader = (
        <Dropzone multiple={false} onDrop={this.handleDrop}>
          <div className="text-center">
            Drop files here <br/>
            Or click here to browse
          </div>
        </Dropzone>
      );
    } else if (this.state.attachments.length > 0) {
      imageUploader = this.state.attachments.map((file, i) => {
        return (
          <div key={i}>
            <img src={file.preview} width="100" />
              <a href="javascript:void(0)" onClick={this.handleImageCancel}>
                <i className="fi-x"></i>
              </a>
          </div>
        );
      });
    } else {
      imageUploader = (
        <a href="javascript:void(0)" onClick={this.handleImageClick}>
          <i className="fi-camera"></i>
        </a>
      );
    }

    return (
      <div className="post-composer">
        <div className="post-composer-field">
          <textarea rows="1" placeholder="Write a new post" onFocus={this.handleFocus} ref="body"></textarea>
        </div>
        <div className="post-composer-actions clearfix">
          <div className="left">
            {imageUploader}
          </div>
          <button className="button right" onClick={this.handleCreate}>Post</button>
        </div>
      </div>
    );
  }
});
