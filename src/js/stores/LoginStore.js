import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _loginState = { error: "", success: false };
let _token = "";
let _user = {username: ""};

// Facebook style store creation.
const LoginStore = assign({}, BaseStore, {
  getLoginState() {
    return _loginState;
  },

  getToken() {
    return _token;
  },

  getUser() {
    return _user;
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
        _loginState.success = true;
        LoginStore.emitChange();
        break;

      case Constants.ActionTypes.LOGIN_FAIL:
        let text = action.text.trim();
        if (text !== '') {
          _loginState.error = text;
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

      case Constants.ActionTypes.FETCHED_USER:
        let user = action.data;
        if (user !== {}) {
          _user = user;
          LoginStore.emitChange();
        }
        break;

    }
  })
});

export default LoginStore;
