import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  getThreads() {
    $.ajax({
      url: Constants.APIROOT + '/messages/threads',
      method: 'GET',
      dataType: 'json'
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_THREADS_SUCCESS,
          data: res
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  }
};
