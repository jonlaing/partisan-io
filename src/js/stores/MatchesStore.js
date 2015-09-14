import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _matches = [];

// Facebook style store creation.
const MatchesStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  getAll() {
    return {
      matches: _matches
    };
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.GET_MATCHES:
        _matches = action.data;
        MatchesStore.emitChange();
        break;

      // add more cases for other actionTypes...
    }
  })
});

export default MatchesStore;
