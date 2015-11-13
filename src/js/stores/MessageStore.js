import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _messages = [];
let _count = 0;

// for some reason, new messages come up twice by default
// gotta make sure they're all unique
function uniqueMessages(messages) {
  var ids = [];

  return messages.filter((msg) => {
    if(ids.indexOf(msg.id) === -1) {
      ids.push(msg.id);
      return true;
    }
  });
}

// Facebook style store creation.
const MessageStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  getMessages() {
    return _messages;
  },

  getCount() {
    return _count;
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.GET_MESSAGES_SUCCESS:
        if (action.data) {
          _messages = uniqueMessages(action.data);
          MessageStore.emitChange();
        }
        break;
      case Constants.ActionTypes.GET_NEW_MESSAGES:
        if (action.data) {
          _messages = uniqueMessages(_messages.concat(action.data));
          MessageStore.emitChange();
        }
        break;
      case Constants.ActionTypes.GET_MESSAGE_COUNT:
        if (action.count !== undefined) {
          _count = action.count;
          MessageStore.emitChange();
        }
        break;

      // add more cases for other actionTypes...
    }
  })
});

export default MessageStore;
