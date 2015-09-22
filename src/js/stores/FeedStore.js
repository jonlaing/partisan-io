import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _feedItems = [];
let _comments = {
  postComments: [],
  postCommentCounts: []
};
let _likes = {
  postLikes: [],
  commentLikes: [],
  liked: []
};

// add private functions to modify data
function _addItems(items) {
  _feedItems = _feedItems.concat(items);
}

function _addItem(item) {
  _feedItems = [item].concat(_feedItems);
}

// add private functions to modify data
function _addCommentCount(data) {
  switch(data.record_type) {
    case "post":
    case "posts":
      for(let i = 0; i < _feedItems.length; i++) {
        if((_feedItems[i].record_type === "post" || _feedItems[i].record_type === "post")
          && data.record_id === _feedItems[i].record_id) {
            _feedItems[i].record.comment_count = data.comment_count;
          }
      }
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
      if(_comments.postComments[id] === undefined) {
        _comments.postComments[id] = [data];
      } else {
        _comments.postComments[id].push(data);
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
      _comments.postComments[data[0].comment.record_id] = data;
      break;
    default:
      break;
  }
}

function _addLike(data) {
  switch(data.record_type) {
    case "post":
    case "posts":
      for(let i = 0; i < _feedItems.length; i++) {
        if((_feedItems[i].record_type === "post" || _feedItems[i].record_type === "post")
          && data.record_id === _feedItems[i].record_id) {
            _feedItems[i].record.like_count = data.like_count;
            _feedItems[i].record.liked = data.liked;
          }
      }
      break;
    case "comment":
    case "comments":
      for(let i = 0; i < _feedItems.length; i++) {
        if((_feedItems[i].record_type === "comment" || _feedItems[i].record_type === "comments")
          && data.record_id === _feedItems[i].record_id) {
            _feedItems[i].record.like_count = data.like_count;
            _feedItems[i].record.liked = data.liked;
          }
      }
      break;
    default:
      break;
  }
}

export default FeedStore;

// Facebook style store creation.
const FeedStore = assign({}, BaseStore, {
  getFeedItems() {
    return _feedItems;
  },

  listComments(type, id) {
    switch(type) {
      case "post":
      case "posts":
        let comments = _comments.postComments[id];
        if(comments === undefined) {
          return [];
        }
        return comments;
      default:
        break;
    }
  },

  countComments(type, id) {
    switch(type) {
      case "post":
      case "posts":
        return _comments.postCommentCounts[id];
      default:
        break;
    }
  },
  //
  // public methods used by Controller-View to operate on data
  getLikes(type, id) {
    switch(type) {
      case "post":
      case "posts":
        return {likeCount: _likes.postLikes[id], liked: _likes.liked[id]};
      case "comment":
      case "comments":
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
        if(action.data) {
          _addItems(action.data);
          FeedStore.emitChange();
        }
        break;
      case Constants.ActionTypes.ADD_FEED_ITEM:
        if(action.data) {
          _addItem(action.data);
          FeedStore.emitChange();
        }
        break;
      // COMMENT ACTIONS
      case Constants.ActionTypes.GET_COMMENT_COUNT_SUCCESS:
        _addCommentCount(action.data);
        FeedStore.emitChange();
        break;
      case Constants.ActionTypes.CREATE_COMMENT_SUCCESS:
        _addComment(action.data);
        FeedStore.emitChange();
        break;
      case Constants.ActionTypes.GET_COMMENTS_SUCCESS:
        _addComments(action.data.comments);
        FeedStore.emitChange();
        break;
      // LIKE ACTIONS
      case Constants.ActionTypes.GET_LIKES_SUCCESS:
        _addLike(action.data);
        FeedStore.emitChange();
        break;
      case Constants.ActionTypes.LIKE_SUCCESS:
        _addLike(action.data);
        FeedStore.emitChange();
        break;

      // add more cases for other actionTypes...
    }
  })
});

export default FeedStore;
