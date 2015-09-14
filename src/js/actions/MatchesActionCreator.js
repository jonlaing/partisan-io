import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  getMatches() {
    $.ajax({
      url: Constants.APIROOT + '/matches',
      method: 'GET',
      dataType: 'json'
    })
      .done(function(res) {
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
