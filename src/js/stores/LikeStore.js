import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// NOTE: I have a lot of switch statements and shit to be ready for other types of records to be liked

// data storage
let _postLikes = [];
let _commentLikes = [];
let _liked = [];

// Facebook style store creation.
const LikeStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  getLikes(type, id) {
    switch(type) {
      case "post":
      case "posts":
        return {likeCount: _postLikes[id], liked: _liked[id]};
      case "comment":
      case "comments":
        return {likeCount: _commentLikes[id], liked: _liked[id]};
      default:
        break;
    }
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.GET_LIKES_SUCCESS:
        _addLike(action.data);
        LikeStore.emitChange();
        break;
      case Constants.ActionTypes.LIKE_SUCCESS:
        _addLike(action.data);
        LikeStore.emitChange();
        break;

    }
  })
});

function _addLike(data) {
  switch(data.record_type) {
    case "post":
    case "posts":
      _postLikes[data.record_id] = data.like_count;
      _liked[data.record_id] = data.liked;
      break;
    case "comment":
    case "comments":
      _commentLikes[data.record_id] = data.like_count;
      _liked[data.record_id] = data.liked;
      break;
    default:
      break;
  }
}

export default LikeStore;
