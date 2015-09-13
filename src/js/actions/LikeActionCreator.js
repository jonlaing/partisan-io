import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  getLikes(type, id) {
    $.ajax({
      url: _likeRoute(type, id),
      method: 'GET',
      dataType: 'json'
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_LIKES_SUCCESS,
          data: res
        });
      })
      .fail(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_LIKES_FAIL,
          data: res
        });
      });
  },

  like(type, id) {
    $.ajax({
      url: _likeRoute(type, id),
      method: 'POST',
      dataType: 'json'
    })
      .done(function(res) {
        console.log(res);
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.LIKE_SUCCESS,
          data: res
        });
      })
      .fail(function(res) {
        console.log(res);
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.LIKE_FAIL,
          data: res
        });
      });
  }

};

function _likeRoute(type, id) {
  return Constants.APIROOT + '/' + type + '/' + id + '/likes';
}
