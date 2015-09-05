/*global $ */
import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  createPost(body) {
    $.ajax({
      url: Constants.APIROOT + '/posts',
      data: {
        body: body
      },
      method: 'POST',
      dataType: 'json'
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.ADD_FEED_ITEM,
          data: res
        });
      })
      .fail(function(res) {
        console.log("fail");
      })
      .always(function(res) {
        console.log(res);
      });
  }
};
