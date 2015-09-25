import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

let _comments = {
  postComments: [],
  postCommentCounts: []
};

function _addComment(data) {
  let id = data.comment.post_id;

  if(_comments.postComments[id] === undefined) {
    _comments.postComments[id] = [data];
  } else {
    _comments.postComments[id].push(data);
  }
}

function _addComments(data) {
  if(data.length < 1) {
    return;
  }

  _comments.postComments[data[0].comment.post_id] = data;
}

function _addCommentLike(data) {
  let id = data.record_id;

  _comments.postComments = _comments.postComments.map((comments) => {
    for(let j = 0; j < comments.length; j++) {
      if(comments[j].comment.id === id) {
        comments[j].like_count = data.like_count;
        comments[j].liked = data.liked;
      }
    }
    return comments;
  });
}

// Facebook style store creation.
const CommentStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  listComments(id) {
    let comments = _comments.postComments[id];
    if(comments === undefined) {
      return [];
    }
    return comments;
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.CREATE_COMMENT_SUCCESS:
        _addComment(action.data);
        CommentStore.emitChange();
        break;
      case Constants.ActionTypes.GET_COMMENTS_SUCCESS:
        _addComments(action.data.comments);
        CommentStore.emitChange();
        break;
      case Constants.ActionTypes.GET_LIKES_SUCCESS:
        if(action.data.record_type === "comment") {
          _addCommentLike(action.data);
          CommentStore.emitChange();
        }
        break;
      case Constants.ActionTypes.LIKE_SUCCESS:
        if(action.data.record_type === "comment") {
          _addCommentLike(action.data);
          CommentStore.emitChange();
        }
        break;

      // add more cases for other actionTypes...
    }
  })
});

export default CommentStore;
