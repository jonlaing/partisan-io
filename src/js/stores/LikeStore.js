import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// NOTE: I have a lot of switch statements and shit to be ready for other types of records to be liked

// data storage
let _postLikes = [];
let _liked = [];

// Facebook style store creation.
const LikeStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  getLikes(type, id) {
    console.log(_postLikes);
    switch(type) {
      case "post":
      case "posts":
        return {likeCount: _postLikes[id], liked: _liked[id]};
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
        console.log(action.data);
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
      console.log(data.record_type);
      console.log(data.like_count);
  switch(data.record_type) {
    case "post":
    case "posts":
      _postLikes[data.record_id] = data.like_count;
      _liked[data.record_id] = data.liked;
      break;
    default:
      break;
  }
}

export default LikeStore;
