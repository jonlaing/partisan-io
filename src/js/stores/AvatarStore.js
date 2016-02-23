import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _avatar = "";

// Facebook style store creation.
const AvatarStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  getAvatar() {
    return _avatar;
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.UPLOAD_AVATAR_SUCCESS:
        if (action.data.avatar_url) {
          _avatar = {
            full: action.data.avatar_url,
            thumb: action.data.avatar_thumbnail_url
          };
          AvatarStore.emitChange();
        }
        break;

      // add more cases for other actionTypes...
    }
  })
});

export default AvatarStore;
