import React from 'react';

import CommentsActionCreator from '../actions/CommentsActionCreator';
import Dropzone from 'react-dropzone';

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
    this.setState({showImageUploader: true});
  },

  handleDrop(files) {
    this.setState({showImageUploader: false, attachments: files});
  },

  handleImageCancel() {
    this.setState({attachments: [], showImageUploader: false});
  },

  render() {
    var uploader;

    if(this.state.showImageUploader === true) {
      uploader = (
        <Dropzone multiple={false} onDrop={this.handleDrop}>
          <div className="text-center">
            Drop files here <br/>
            Or click here to browse
          </div>
        </Dropzone>
      );
    } else if (this.state.attachments.length > 0) {
      uploader = this.state.attachments.map((file, i) => {
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
      uploader = (
        <a href="javascript:void(0)" onClick={this.handlePhoto}><i className="fi-camera"></i></a>
      );
    }

    return (
      <div className="comment-composer">
        <div className="row collapse">
          <div className="large-10 columns">
            <textarea type="text" placeholder="Type your comment here..." ref="comment" ></textarea>
          </div>
          <div className="large-2 columns">
            <button onClick={this.handleSubmit}>Comment</button>
          </div>
        </div>
        <div className="row">
          {uploader}
        </div>
      </div>
    );
  }
});
