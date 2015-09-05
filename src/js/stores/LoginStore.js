import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _data = { error: "", success: false };
let _token = "";

// Facebook style store creation.
const LoginStore = assign({}, BaseStore, {
  getState() {
    return _data;
  },

  getToken() {
    return _token;
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.LOGIN_SUCCESS:
        let token = action.data.token;
        if(token !== '') {
          _token = token;
        }
        _data.success = true;
        LoginStore.emitChange();
        break;

      case Constants.ActionTypes.LOGIN_FAIL:
        let text = action.text.trim();
        if (text !== '') {
          _data.error = text;
          LoginStore.emitChange();
        }
        break;

      // this is for intentionally logging out
      case Constants.ActionTypes.LOGOUT:
        window.location = "/login.html";
        break;

      // this is incase your session expired
      case Constants.ActionTypes.LOGGED_OUT:
        window.location = "/login.html";
        break;

    }
  })
});

export default LoginStore;
