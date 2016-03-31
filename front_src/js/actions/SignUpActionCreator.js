import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  signUp(user) {
    $.ajax({
      url: Constants.APIROOT + '/users',
      data: {
        email: user.email,
        username: user.username,
        postal_code: user.postalCode,
        password: user.password,
        password_confirm: user.passwordConfirm
      },
      method: 'POST',
      dataType: 'json'
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.SIGN_UP_SUCCESS,
          data: res
        });
      })
      .fail(function(res) {
        let errors = JSON.parse(res.responseText);
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.SIGN_UP_FAIL,
          errors: errors
        });
      });
  },

  checkUnique(username) {
    if(username.length > 0) {
      $.ajax({
        url: Constants.APIROOT + '/user/check_unique',
        data: { username: username },
        method: 'GET',
        dataType: 'json'
      })
        .done(function() {
          Dispatcher.handleViewAction({
            type: Constants.ActionTypes.USERNAME_UNIQUE
          });
        })
        .fail(function() {
          Dispatcher.handleViewAction({
            type: Constants.ActionTypes.USERNAME_NOT_UNIQUE
          });
        });
    } else {
      Dispatcher.handleViewAction({
        type: Constants.ActionTypes.USERNAME_BLANK
      });
    }

  }
};
