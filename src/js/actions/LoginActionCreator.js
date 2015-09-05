/*global $ */
import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  login(email, password) {
    $.ajax({
      url: Constants.APIROOT + '/login',
      data: {
        email: email,
        password: password
      },
      method: 'POST',
      dataType: 'json'
    })
      .done(function(data) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.LOGIN_SUCCESS,
          data: data
        });
      })
      .fail(function() {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.LOGIN_FAIL,
          text: "Login failed."
        });
      });
  },

  logout() {
    $.get(Constants.APIROOT + '/logout')
      .always(function() {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.LOGOUT
        });
      });
  }
};
