import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  getAll() {
    $.ajax({
      url: Constants.APIROOT + '/friendships/',
      method: 'GET',
      dataType: 'json'
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_FRIENDSHIPS_SUCCESS,
          data: res.friendships
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  },

  getFriendship(friendID) {
    $.ajax({
      url: Constants.APIROOT + '/friendships/' + friendID,
      method: 'GET',
      dataType: 'json'
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_FRIENDSHIP_SUCCESS,
          data: {id: friendID, friendship: res}
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  },

  requestFriendship(friendID) {
    $.ajax({
      url: Constants.APIROOT + '/friendships',
      data: { "friend_id": friendID },
      method: 'POST',
      dataType: 'json'
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.REQUEST_FRIENDSHIP_SUCCESS,
          data: {id: friendID, friendship: res}
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  },

  confirmFriendship(friendID) {
    $.ajax({
      url: Constants.APIROOT + '/friendships',
      data: { "friend_id": friendID },
      method: 'PATCH',
      dataType: 'json'
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.CONFIRM_FRIENDSHIP_SUCCESS,
          data: {id: friendID, friendship: res}
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  }
};
