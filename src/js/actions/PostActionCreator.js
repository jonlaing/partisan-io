/*global $, FormData */
import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  createPost(body, attachments) {
    var request = new FormData();

    attachments.forEach(function(value) {
      request.append('attachment', value);
    });

    request.append('body', body);

    $.ajax({
      url: Constants.APIROOT + '/posts',
      data: request,
      cache: false,
      method: 'POST',
      processData: false,
      contentType: false
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.ADD_FEED_ITEM,
          data: res
        });
      })
      .fail(function(res) {
        console.log(res);
      })
      .always(function(res) {
        console.log(res);
      });
  }
};
