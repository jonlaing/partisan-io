/*global WebSocket */
import Dispatcher from '../Dispatcher';
import Constants from '../Constants';


export default {
  getNotificationList() {
    $.ajax({
      url: Constants.APIROOT + '/notifications',
      method: 'GET',
      dataType: 'json'
    })
      .done(function(res) {
        console.log(res);
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_NOTIFICATIONS_SUCCESS,
          data: res
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  },

  getNotificationCount() {
    // TODO: uncomment for prod
    // var _socket = new WebSocket("ws://localhost:4000" + Constants.APIROOT + "/notifications/count");

    // _socket.onmessage = (res) => {
    //   let data = JSON.parse(res.data);
    //   Dispatcher.handleViewAction({
    //     type: Constants.ActionTypes.GET_NOTIFICATION_COUNT,
    //     data: data.count
    //   });
    // };

    // _socket.onopen = () => {
    //   window.setInterval(() => {
    //     _socket.send(0);
    //   }, 5000);
    // };
  }
};
