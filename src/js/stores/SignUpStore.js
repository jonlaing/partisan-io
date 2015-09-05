import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _errors = [];
let _success = false;

// Facebook style store creation.
const SignUpStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  getState() {
    return {
      errors: _errors,
      success: _success
    };
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.SIGN_UP_SUCCESS:
        _success = true;
        SignUpStore.emitChange();
        break;

      case Constants.ActionTypes.SIGN_UP_ERROR:
        _errors = action.errors;
        SignUpStore.emitChange();
        break;
    }
  })
});

export default SignUpStore;
