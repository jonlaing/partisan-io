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
        // Logged out
        if(res.status === 401) {
          Dispatcher.handleViewAction({
            type: Constants.ActionTypes.LOGGED_OUT
          });
        } else if (res.status === 404) {
          Dispatcher.handleViewAction({
            type: Constants.ActionTypes.NO_FRIENDS
          });
        }
      });
  },
  getByUser(userID) {
    $.ajax({
      url: '/api/v1/feed/' + userID,
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
        // Logged out
        if(res.status === 401) {
          Dispatcher.handleViewAction({
            type: Constants.ActionTypes.LOGGED_OUT
          });
        }
      });
  },
  addItem(data) {
    Dispatcher.handleViewAction({
      type: Constants.ActionTypes.ADD_FEED_ITEM,
      data: data
    });
  }
};
