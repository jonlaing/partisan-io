import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  updateLocation(postalCode) {
    $.ajax({
      url: Constants.APIROOT + '/users',
      data: {
        "postal_code": postalCode
      },
      method: 'PATCH',
      dataType: 'json'
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.UPDATE_USER_SUCCESS,
          user: res
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  },

  updateGender(gender) {
    $.ajax({
      url: Constants.APIROOT + '/users',
      data: {
        "gender": gender
      },
      method: 'PATCH',
      dataType: 'json'
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.UPDATE_USER_SUCCESS,
          user: res
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  },

  updateLookingFor(val) {
    $.ajax({
      url: Constants.APIROOT + '/profile',
      data: {
        "looking_for": val
      },
      method: 'PATCH',
      dataType: 'json'
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.UPDATE_PROFILE_SUCCESS,
          profile: res
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  }
};
