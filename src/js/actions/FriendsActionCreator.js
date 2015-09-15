import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
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
  }
};
