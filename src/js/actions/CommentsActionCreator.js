/*global FormData */
import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  getList(type, id) {
    $.ajax({
      url: _commentRoute(id),
      method: 'GET',
      dataMethod: 'json'
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_COMMENTS_SUCCESS,
          data: res
        });
      })
      .fail(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_COMMENTS_FAIL,
          data: res
        });
      });
  },

  getCount(type, id) {
    $.ajax({
      url: _commentRoute(id) + '/count',
      method: 'GET',
      dataMethod: 'json'
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_COMMENT_COUNT_SUCCESS,
          data: res
        });
      })
      .fail(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_COMMENT_COUNT_FAIL,
          data: res
        });
      });
  },

  create(postID, body, attachments) {
    var request = new FormData();

    attachments.forEach(function(value) {
      request.append('attachment', value);
    });

    request.append('post_id', postID);
    request.append('body', body);

    $.ajax({
      url: Constants.APIROOT + '/comments',
      data: request,
      cache: false,
      method: 'POST',
      processData: false,
      contentType: false
    })
      .done(function(res) {
        let data = {
          post_id: res.comment.post_id,
          comment_count: res.count
        };

        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.CREATE_COMMENT_SUCCESS,
          data: res
        });
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_COMMENT_COUNT_SUCCESS,
          data: data
        });
      })
      .fail(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.CREATE_COMMENT_FAIL,
          data: res
        });
      });
  }

};

function _commentRoute(id) {
  return Constants.APIROOT + '/posts/' + id + '/comments';
}
