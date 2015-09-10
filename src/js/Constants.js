import keyMirror from 'react/lib/keyMirror';

export default {
  APIROOT: '/api/v1',

  // event name triggered from store, listened to by views
  CHANGE_EVENT: 'change',

  // Each time you add an action, add it here
  ActionTypes: keyMirror({
    GET_FEED: null,
    ADD_FEED_ITEM: null,

    LOGIN_SUCCESS: null,
    LOGIN_FAIL: null,
    LOGOUT: null,
    LOGGED_OUT: null,
    FETCHED_USER: null,

    SIGN_UP_SUCCESS: null,
    SIGN_UP_FAIL: null,

    GET_QUESTION_SUCESS: null,
    GET_QUESTION_FAIL: null
  }),

  ActionSources: keyMirror({
    SERVER_ACTION: null,
    VIEW_ACTION: null
  })
};
