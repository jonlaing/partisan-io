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
          data: res.notifications
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  },

  getNotificationCount() {
    var domain;
    let url = window.location.href;

    //find & remove protocol (http, ftp, etc.) and get domain
    if (url.indexOf("://") > -1) {
        domain = url.split('/')[2];
    }
    else {
        domain = url.split('/')[0];
    }

    var _socket = new WebSocket("wss://" + domain + Constants.APIROOT + "/notifications/count");

    _socket.onmessage = (res) => {
      let data = JSON.parse(res.data);
      Dispatcher.handleViewAction({
        type: Constants.ActionTypes.GET_NOTIFICATION_COUNT,
        data: data.count
      });
    };

    _socket.onopen = () => {
      _socket.send(0); // grab it once off the bat
      window.setInterval(() => {
        _socket.send(0);
      }, 5000);
    };

  }
};
