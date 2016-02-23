import React from 'react';

import Icon from 'react-fontawesome';

import CommentsActionCreator from '../actions/CommentsActionCreator';

export default React.createClass({
  getInitialState() {
    return { showImageUploader: false, attachments: [] };
  },

  handleSubmit() {
    let body = $(React.findDOMNode(this.refs.comment));
    CommentsActionCreator.create(this.props.id, body.val(), this.state.attachments);
    body.val('');
  },

  handlePhoto() {
    // this.setState({showImageUploader: true});
    $(React.findDOMNode(this.refs.file)).click();
  },

  handleDrop(files) {
    this.setState({showImageUploader: false, attachments: files});
  },

  handleImageCancel() {
    this.setState({attachments: [], showImageUploader: false});
    React.findDOMNode(this.refs.file).value = null;
  },

  handleFileChange(e) {
    let attachments = e.target.files[0];
    this.setState({attachments: [attachments] });
  },

  render() {
    var imagePreview;

    if (this.state.attachments.length > 0) {
      imagePreview = this.state.attachments.map((file, i) => {
        return (
          <div key={i} className="comment-composer-image">
            <img src={window.URL.createObjectURL(file)} width="100" />
            <a href="javascript:void(0)" onClick={this.handleImageCancel} className="comment-composer-upload-cancel">
              <Icon name="times" />
            </a>
          </div>
        );
      });
    } else {
      imagePreview = "";
    }

    return (
      <div className="comment-composer">
        <div className="comment-composer-form">
          <div className="comment-composer-input">
            <textarea type="text" placeholder="Type your comment here..." ref="comment" ></textarea>
            <label htmlFor={"comment-file-" + this.props.id} className="comment-composer-uploader">
              <a href="javascript:void(0)" onClick={this.handlePhoto}>
                <Icon name="camera-retro" />
              </a>
              <input
                type="file"
                name={"comment-file-" + this.props.id} style={{display: 'none'}}
                onChange={this.handleFileChange}
                accept="image/jpeg,image/png"
                ref="file" />
            </label>
          </div>
          <button onClick={this.handleSubmit}>Comment</button>
        </div>
        <div className="row">
          {imagePreview}
        </div>
      </div>
    );
  }
});
