import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _suggestions = [];

// Facebook style store creation.
const PostComposerStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  getUserSuggestions() {
    return {
      usernameSuggestions: _suggestions
    };
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.GET_USERNAME_SUGGESTIONS:
        if (action.suggestions !== null) {
          _suggestions = action.suggestions;
        } else {
          _suggestions = [];
        }
        PostComposerStore.emitChange();
        break;

      // add more cases for other actionTypes...
    }
  })
});

export default PostComposerStore;
