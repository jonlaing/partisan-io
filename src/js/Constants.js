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
    USERNAME_UNIQUE: null,
    USERNAME_NOT_UNIQUE: null,
    USERNAME_BLANK: null,

    GET_QUESTION_SUCESS: null,
    GET_QUESTION_FAIL: null,

    LIKE_SUCCESS: null,
    LIKE_FAIL: null,
    GET_LIKES_SUCCESS: null,
    GET_LIKES_FAIL: null,
    UNLIKE_SUCCESS: null,
    UNLIKE_FAIL: null,

    GET_COMMENT_COUNT_SUCCESS: null,
    GET_COMMENT_COUNT_FAIL: null,
    GET_COMMENTS_SUCCESS: null,
    GET_COMMENTS_FAIL: null,
    CREATE_COMMENT_SUCCESS: null,
    CREATE_COMMENT_FAIL: null,

    GET_MATCHES: null,

    GET_FRIENDSHIP_SUCCESS: null,
    REQUEST_FRIENDSHIP_SUCCESS: null,
    CONFIRM_FRIENDSHIP_SUCCESS: null,

    UPLOAD_AVATAR_SUCCESS: null,

    UPDATE_USER_SUCCESS: null,
    UPDATE_PROFILE_SUCCESS: null,

    GET_NOTIFICATIONS_SUCCESS: null,
    GET_NOTIFICATION_COUNT: null,

    BEGIN_FLAG: null,
    CANCEL_FLAG: null,
    SUBMIT_FLAG: null,

    GET_USERNAME_SUGGESTIONS: null
  }),

  ActionSources: keyMirror({
    SERVER_ACTION: null,
    VIEW_ACTION: null
  })
};
