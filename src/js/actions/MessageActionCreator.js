/*global WebSocket */
import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

import moment from 'moment';

export default {
  getMessages(threadID) {
    $.ajax({
      url: Constants.APIROOT + '/messages/threads/' + threadID,
      method: 'GET',
      dataType: 'json'
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_MESSAGES_SUCCESS,
          data: res.messages
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  },
  messageSocket(threadID) {
    var domain;
    let url = window.location.href;

    //find & remove protocol (http, ftp, etc.) and get domain
    if (url.indexOf("://") > -1) {
        domain = url.split('/')[2];
    }
    else {
        domain = url.split('/')[0];
    }

    var _socket, sendInterval, reopenInterval;

    var start = function() {
      if(!_socket) {
        _socket = new WebSocket("ws://" + domain + Constants.APIROOT + "/messages/threads/" + threadID + "/socket");
      } else {
        return;
      }

      _socket.onmessage = (res) => {
        let data = JSON.parse(res.data);
        if(data.messages) {
          Dispatcher.handleViewAction({
            type: Constants.ActionTypes.GET_NEW_MESSAGES,
            data: data.messages
          });
        }
      };

      _socket.onopen = () => {
        window.clearInterval(reopenInterval);

        let lastNow = moment(Date.now()).unix();

        sendInterval = window.setInterval(() => {
          if(!_socket || _socket.readyState === 2 || _socket.readyState === 3) {
            window.clearInterval(sendInterval);
            return;
          }

          _socket.send(lastNow.toString());
          lastNow = moment(Date.now()).unix();
        }, 500);
      };

      _socket.onclose = () => {
        reopenInterval = window.setInterval(() => {
          if(!_socket || _socket.readyState === 0 || _socket.readyState === 1) {
            window.clearInterval(reopenInterval);
            return;
          }

          _socket = null;
          start();
        }, 5000);
      };
    };

    start();

  },
  sendMessage(threadID, text) {
    $.ajax({
      url: Constants.APIROOT + '/messages/threads/' + threadID,
      data: {
        body: text
      },
      method: 'POST',
      dataType: 'json'
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.SEND_MESSAGE_SUCCESS,
          data: res
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  }
};
