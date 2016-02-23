import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _hits = [];

function addHits(hits) {
  _hits = hits;
}

function _addLike(data) {
  let id = data.record_id;
  switch(data.record_type) {
    case "post":
      _hits[id].like_count = data.like_count;
      _hits[id].liked = data.liked;
      break;
    default:
      break;
  }
}

// Facebook style store creation.
const HashtagStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  getAll() {
    return {
      hits: _hits
    };
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.SEARCH_HASHTAG_SUCCESS:
        let data = action.data;

      if (data.length !== 0) {
          addHits(data);
          HashtagStore.emitChange();
        }
        break;
      case Constants.ActionTypes.LIKE_SUCCESS:
        _addLike(action.data);
        HashtagStore.emitChange();
        break;

      // add more cases for other actionTypes...
    }
  })
});

export default HashtagStore;
