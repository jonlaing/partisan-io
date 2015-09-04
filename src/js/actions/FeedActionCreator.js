/*global $ */
import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  getFeed() {
    $.ajax({
      url: '/api/v1/feed',
      dataType: 'json',
      method: 'GET'
    })
      .done(function(res) {
        let data = res.feed_items;

        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_FEED,
          data: data
        });
      })
      .fail(function(res) {
        console.log('something went wrong');
        console.log(res);
      });
  },
  addItem(data) {
    Dispatcher.handleViewAction({
      type: Constants.ActionTypes.ADD_FEED_ITEM,
      data: data
    });
  }
};
