import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  search(query) {
    var encoded = encodeURIComponent(query);
    $.ajax({
      url: Constants.APIROOT + '/hashtags',
      data: { q: encoded },
      method: 'GET',
      dataType: 'json'
    })
      .done(function(res) {
        console.log(res);
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.ADD_TASK,
          data: res
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  }
};
