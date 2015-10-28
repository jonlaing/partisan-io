import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _friends = []; // holds data for "Friender" button
let _friendships = []; // holds data for Friend List

// add private functions to modify data
function addFriend(id, friendship) {
  _friends[id] = friendship;
}

// Facebook style store creation.
const FriendsStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  getAll() {
    return { friendships: _friendships };
  },

  getFriend(id) {
    let friendship = _friends[id];
    if(friendship !== undefined) {
      return friendship;
    } else {
      return {};
    }
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.GET_FRIENDSHIPS_SUCCESS:
        if(action.data) {
          _friendships = action.data;
          FriendsStore.emitChange();
        }
        break;
      case Constants.ActionTypes.GET_FRIENDSHIP_SUCCESS:
        if(action.data.friendship) {
          addFriend(action.data.id, action.data.friendship);
          FriendsStore.emitChange();
        }
        break;
      case Constants.ActionTypes.REQUEST_FRIENDSHIP_SUCCESS:
        if (action.data.friendship) {
          addFriend(action.data.id, action.data.friendship);
          FriendsStore.emitChange();
        }
        break;
      case Constants.ActionTypes.CONFIRM_FRIENDSHIP_SUCCESS:
        if (action.data.friendship) {
          // will replace the old friendship
          addFriend(action.data.id, action.data.friendship);
          FriendsStore.emitChange();
        }
        break;

      // add more cases for other actionTypes...
    }
  })
});

export default FriendsStore;
