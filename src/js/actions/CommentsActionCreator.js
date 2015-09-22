import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  getList(type, id) {
    $.ajax({
      url: _commentRoute(type, id),
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
      url: _commentRoute(type, id) + '/count',
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

  create(comment) {
    $.ajax({
      url: Constants.APIROOT + '/comments',
      data: JSON.stringify(comment),
      method: 'POST',
      dataMethod: 'json'
    })
      .done(function(res) {
        let data = {
          record_type: res.comment.record_type,
          record_id: res.comment.record_id,
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

function _commentRoute(type, id) {
  return Constants.APIROOT + '/' + type + '/' + id + '/comments';
}
