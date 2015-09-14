import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _postCommentsCounts = [];
let _postComments = [];

// add private functions to modify data
function _addCommentCount(data) {
  switch(data.record_type) {
    case "post":
    case "posts":
      _postCommentsCounts[data.record_id] = data.comment_count;
      break;
    default:
      break;
  }
}

function _addComment(data) {
  let id = data.comment.record_id;
  switch(data.comment.record_type) {
    case "post":
    case "posts":
      if(_postComments[id] === undefined) {
        _postComments[id] = [data];
      } else {
        _postComments[id].push(data);
      }
      break;
    default:
      break;
  }
}

function _addComments(data) {
  if(data.length < 1) {
    return;
  }

  switch(data[0].comment.record_type) {
    case "post":
    case "posts":
      _postComments[data[0].comment.record_id] = data;
      break;
    default:
      break;
  }
}

// Facebook style store creation.
const CommentsStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  list(type, id) {
    switch(type) {
      case "post":
      case "posts":
        let comments = _postComments[id];
        if(comments === undefined) {
          return [];
        }
        return comments;
      default:
        break;
    }
  },

  count(type, id) {
    switch(type) {
      case "post":
      case "posts":
        return _postCommentsCounts[id];
      default:
        break;
    }
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.GET_COMMENT_COUNT_SUCCESS:
        _addCommentCount(action.data);
        CommentsStore.emitChange();
        break;
      case Constants.ActionTypes.CREATE_COMMENT_SUCCESS:
        _addComment(action.data);
        CommentsStore.emitChange();
        break;
      case Constants.ActionTypes.GET_COMMENTS_SUCCESS:
        _addComments(action.data.comments);
        CommentsStore.emitChange();
        break;
    }
  })
});

export default CommentsStore;
