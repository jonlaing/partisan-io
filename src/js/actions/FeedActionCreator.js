/*global $ */
import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  getFeed() {
    $.ajax({
      url: '/api/v1/feed',
      dataType: 'jsonp',
      method: 'GET'
    })
      .always(function(res) {
        let data = JSON.parse(res.responseText);

        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_FEED,
          data: data
        });
      });
  }
};
