/*global $ */
import React from 'react';
import Dropzone from 'react-dropzone';

import Icon from 'react-fontawesome';

import PostActionCreator from '../actions/PostActionCreator';
import PostComposerStore from '../stores/PostComposerStore';

const TAB = 9;
const UP = 38;
const DOWN = 40;

export default React.createClass({
  getInitialState() {
    return { showImageUploader: false, attachments: [], usernameSuggestions: [], usernameIndex: -1 };
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

  handleChange(e) {
    let text = e.target.value;
    let matches = text.match(/(?:\s|^)@([a-zA-Z0-9_]+)$/);

    if(matches !== null) {
      PostActionCreator.suggestUsers(matches[1]);
    } else {
      this.setState({currentUserTag: "", usernameSuggestions: [], usernameIndex: -1});
    }
  },

  handleKeyDown(e) {
    var newBody, index;

    if(this.state.usernameSuggestions.length > 0) {
      switch(e.keyCode) {
      case TAB:
        e.preventDefault();
        newBody = e.target.value.replace(/@[a-zA-Z0-9_]+$/, "@" + this.state.usernameSuggestions[this.state.usernameIndex]);
        $(e.target).val(newBody);
        break;
      case DOWN:
        e.preventDefault();
        index = this.state.usernameIndex + 1;

        console.log(index);

        if(index < this.state.usernameSuggestions.length) {
          this.setState({usernameIndex: index});
          newBody = e.target.value.replace(/@[a-zA-Z0-9_]+$/, "@" + this.state.usernameSuggestions[index]);
          $(e.target).val(newBody);
        }
        break;
      case UP:
        e.preventDefault();
        index = this.state.usernameIndex - 1;

        console.log(index);

        if(index >= 0) {
          this.setState({usernameIndex: index});
          newBody = e.target.value.replace(/@[a-zA-Z0-9_]+$/, "@" + this.state.usernameSuggestions[index]);
          $(e.target).val(newBody);
        }
        break;
      default:
        break;
      }
    }
  },

  handleSuggestionClick(e) {
    console.log(e.target.innerHTML);
  },

  componentDidMount() {
    PostComposerStore.addChangeListener(this._onChange);
  },

  componentWillUnmount() {
    PostComposerStore.removeChangeListener(this._onChange);
  },

  render() {
    var imageUploader, usernameList;

    if(this.state.showImageUploader === true) {
      imageUploader = (
        <Dropzone multiple={false} onDrop={this.handleDrop} className="post-dropzone" activeClassName="post-dropzone-active">
          <Icon name="download" />
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
        <a href="javascript:void(0)" onClick={this.handleImageClick} className="post-type">
          <Icon name="camera-retro" />
        </a>
      );
    }

    usernameList = this.state.usernameSuggestions.map((suggestion) => <li key={suggestion} onClick={this.handleSuggestionClick}> {suggestion} </li> );

    return (
      <div className="post-composer">
        <div className="post-composer-field">
          <textarea rows="1" placeholder="Write a new post" onFocus={this.handleFocus} onKeyDown={this.handleKeyDown} onChange={this.handleChange} ref="body"></textarea>
        </div>
        <div className="post-composer-actions">
          <div>
            {imageUploader}
          </div>
          <button className="button" onClick={this.handleCreate}>Post</button>
        </div>
        <div className="post-composer-usernames">
          <ul>
            {usernameList}
          </ul>
        </div>
      </div>
    );
  },

  _onChange() {
    let state = PostComposerStore.getUserSuggestions();
    this.setState(state);
  }
});
