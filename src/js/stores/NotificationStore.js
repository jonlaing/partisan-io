import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _notifications = [];
let _notificationCount = 0;

// Facebook style store creation.
const NotificationStore = assign({}, BaseStore, {
  getCount() {
    return { count: _notificationCount };
  },

  getAll() {
    return { count: _notificationCount, notifications: _notifications };
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.GET_NOTIFICATIONS_SUCCESS:
        let notifs = action.data;

        if (notifs.length > 0) {
          _notifications = notifs;
          NotificationStore.emitChange();
        }
        break;
      case Constants.ActionTypes.GET_NOTIFICATION_COUNT:
        let count = action.data;

        if (count !== _notificationCount) {
          // TODO: Make this better!
          if(count > 0) {
            $('title').text('My Feed (' + count + ')');
          } else {
            $('title').text('My Feed');
          }
          _notificationCount = count;
          NotificationStore.emitChange();
        }
        break;

    }
  })
});

export default NotificationStore;
