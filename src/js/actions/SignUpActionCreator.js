import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  signUp(user) {
    $.ajax({
      url: Constants.APIROOT + '/users',
      data: {
        email: user.email,
        username: user.username,
        full_name: user.fullName,
        postal_code: user.postalCode,
        password: user.password,
        password_confirm: user.passwordConfirm
      },
      method: 'POST',
      dataType: 'json'
    })
      .done(function(res) {
        console.log(res);
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
  }
};
