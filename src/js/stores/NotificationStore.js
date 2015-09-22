import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

import FeedStore from './FeedStore';

// data storage
let _notifications = [];
let _notificationCount = 0;

// Facebook style store creation.
const NotificationStore = assign({}, BaseStore, {
  getCount() {
    return { count: _notificationCount };
  },

  getAll() {
    return { count: this.getCount(), notifications: _notifications };
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
          $('title').text('My Feed (' + count + ')');
          _notificationCount = count;
          NotificationStore.emitChange();
        }
        break;

    }
  })
});

export default NotificationStore;
