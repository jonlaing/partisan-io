/*global $ */
import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  signUp(user) {
    $.ajax({
      url: Constants.APIROOT + '/users',
      data: {
        email: user.email,
        username: user.username,
        full_name: user.full_name,
        password: user.password,
        password_confirm: user.password_confirm
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
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.SIGN_UP_FAIL,
          data: res
        });
      });
  }
};
