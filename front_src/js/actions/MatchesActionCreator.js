import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  getMatches(distance, gender, minAge, maxAge) {
    $.ajax({
      url: Constants.APIROOT + '/matches',
      data: {
        distance: distance,
        gender: gender,
        minAge: minAge,
        maxAge: maxAge
      },
      method: 'GET',
      dataType: 'json'
    })
      .done(function(res) {
        console.log(res);
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_MATCHES,
          data: res
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  }
};
