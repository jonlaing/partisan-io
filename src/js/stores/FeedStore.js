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
  for(let i = 0; i < _feedItems.length; i++) {
    if(data.post_id === _feedItems[i].record_id) {
      _feedItems[i].record.comment_count = data.comment_count;
    }
  }
}

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

function _addLike(data) {
  let id = data.record_id;
  switch(data.record_type) {
    case "post":
      for(let i = 0; i < _feedItems.length; i++) {
        if(data.record_id === _feedItems[i].record_id) {
            _feedItems[i].record.like_count = data.like_count;
            _feedItems[i].record.liked = data.liked;
          }
      }
      break;
    case "comment":
      _comments.postComments = _comments.postComments.map((comments) => {
        for(let j = 0; j < comments.length; j++) {
          if(comments[j].comment.id === id) {
            comments[j].like_count = data.like_count;
            comments[j].liked = data.liked;
          }
        }
        return comments;
      });
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

  listComments(id) {
    let comments = _comments.postComments[id];
    if(comments === undefined) {
      return [];
    }
    return comments;
  },

  countComments(id) {
      return _comments.postCommentCounts[id];
  },
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
