import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _feedItems = [];

let _noFriends = false;

let _likes = {
  postLikes: [],
  commentLikes: [],
  liked: []
};

let _modals = {
  flag: {
    show: false,
    id: 0,
    type: ""
  }
};

// add private functions to modify data
function _addItems(items) {
  _feedItems = _feedItems.concat(items);
}

function _prependItems(items) {
  _feedItems = items.concat(_feedItems);
}

function _addItem(item) {
  _feedItems = [item].concat(_feedItems);
}

// add private functions to modify data
function _addCommentCount(data) {
  for(let i = 0; i < _feedItems.length; i++) {
    if(data.post_id === _feedItems[i].record_id) {
      _feedItems[i].record.comment_count = data.comment_count;
    }
  }
}

function _addPostLike(data) {
  for(let i = 0; i < _feedItems.length; i++) {
    if(data.record_id === _feedItems[i].record_id) {
        _feedItems[i].record.like_count = data.like_count;
        _feedItems[i].record.liked = data.liked;
      }
  }
}

export default FeedStore;

// Facebook style store creation.
const FeedStore = assign({}, BaseStore, {
  getState() {
    return { feed: _feedItems, modals: _modals, noFriends: _noFriends, scrollLoading: false };
  },

  // countComments(id) {
  //     return _comments.postCommentCounts[id];
  // },
  //
  // public methods used by Controller-View to operate on data
  getLikes(type, id) {
    switch(type) {
      case "post":
        return {likeCount: _likes.postLikes[id], liked: _likes.liked[id]};
      case "comment":
        return {likeCount: _likes.commentLikes[id], liked: _likes.liked[id]};
      default:
        break;
    }
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      // FEED ACTIONS
      case Constants.ActionTypes.GET_FEED:
      case Constants.ActionTypes.GET_FEED_PAGE:
        if(action.data) {
          _addItems(action.data);
          FeedStore.emitChange();
        }
        break;
      case Constants.ActionTypes.GET_NEW_FEED_ITEMS:
        if(action.data.length > 0) {
          console.log(action.data);
          _prependItems(action.data);
          FeedStore.emitChange();
        }
        break;
      case Constants.ActionTypes.NO_FRIENDS:
        _noFriends = true;
        FeedStore.emitChange();
        break;
      case Constants.ActionTypes.ADD_FEED_ITEM:
        if(action.data) {
          _addItem(action.data);
          FeedStore.emitChange();
        }
        break;
      case Constants.ActionTypes.GET_COMMENT_COUNT_SUCCESS:
        _addCommentCount(action.data);
        FeedStore.emitChange();
        break;
      // LIKE ACTIONS
      case Constants.ActionTypes.GET_LIKES_SUCCESS:
        // we only like for posts on the feed
        // comment posts are taken care of in CommentStore
        if(action.data.record_type === "post") {
          _addPostLike(action.data);
          FeedStore.emitChange();
        }
        break;
      case Constants.ActionTypes.LIKE_SUCCESS:
        if(action.data.record_type === "post") {
          _addPostLike(action.data);
          FeedStore.emitChange();
        }
        break;
      // FLAG ACTIONS
      case Constants.ActionTypes.BEGIN_FLAG:
        _modals.flag.show = true;
        _modals.flag.id = action.id;
        _modals.flag.type = action.recordType;
        FeedStore.emitChange();
        break;
      case Constants.ActionTypes.CANCEL_FLAG:
      case Constants.ActionTypes.SUBMIT_FLAG:
        _modals.flag.show = false;
        _modals.flag.id = 0;
        _modals.flag.type = "";
        FeedStore.emitChange();
        break;

      // add more cases for other actionTypes...
    }
  })
});

export default FeedStore;
