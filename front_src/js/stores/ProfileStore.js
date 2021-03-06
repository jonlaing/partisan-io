import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _user = {};
let _profile = {};

// Facebook style store creation.
const ProfileStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  getProfile() {
    var state = {};
    if(_user.id !== undefined) {
      state.user = _user;
    }

    if(_profile.id !== undefined) {
      state.profile = _profile;
    }

    return state;
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.UPDATE_USER_SUCCESS:
        let user = action.user;

        if (user !== undefined) {
          _user = user;
          ProfileStore.emitChange();
        }
        break;

      case Constants.ActionTypes.UPDATE_PROFILE_SUCCESS:
        let profile = action.profile;

        if (profile !== undefined) {
          _profile = profile;
          ProfileStore.emitChange();
        }
        break;

      // add more cases for other actionTypes...
    }
  })
});

export default ProfileStore;
