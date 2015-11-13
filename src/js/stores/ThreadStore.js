import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

let _threads = [];
let _inactive = [];
let _currentThread = 0;

// Facebook style store creation.
const ThreadStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  getList() {
      return _threads;
  },

  getInactive() {
    return _inactive;
  },

  getCurrentThread() {
    return _currentThread;
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.GET_THREADS_SUCCESS:
        if (action.data) {
          if(action.data.threads !== null) {
            _threads = action.data.threads.filter((thread, i, arr) => {
              return arr.indexOf(thread) === i;
            });
          }

          if(action.data.inactive !== null) {
            _inactive = action.data.inactive.filter((user, i, arr) => {
              return arr.indexOf(user) === i;
            });
          }

          ThreadStore.emitChange();
        }
        break;
      case Constants.ActionTypes.CREATE_THREAD_SUCCESS:
        if(action.thread) {
          _inactive = _inactive.filter((user) => {
            return user.id !== action.thread.thread_user.user_id;
          });

          _threads = _threads.push(action.thread);

          _currentThread = action.thread.thread_user.thread_id;
          ThreadStore.emitChange();
        }
        break;
      case Constants.ActionTypes.SWITCH_THREADS:
        if (action.threadID) {
          _currentThread = action.threadID;
          _threads = _threads.map((thread) => {
            if(thread.thread_user.thread_id === action.threadID) {
              thread.has_unread = false;
            }
            return thread;
          });
          ThreadStore.emitChange();
        }
        break;

      // add more cases for other actionTypes...
    }
  })
});

export default ThreadStore;
