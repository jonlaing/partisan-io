import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _feedItems = [];

// add private functions to modify data
function _addItems(items) {
  _feedItems = _feedItems.concat(items);
}

// Facebook style store creation.
const FeedStore = assign({}, BaseStore, {
  getAll() {
    return _feedItems;
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.GET_FEED:
        let data = action.data;
        if(data) {
          _addItems(data);
          FeedStore.emitChange();
        }
        break;

      // add more cases for other actionTypes...
    }
  })
});

export default FeedStore;
