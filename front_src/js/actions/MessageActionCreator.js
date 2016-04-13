/*global WebSocket */
import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

import moment from 'moment';

var _messageSocket;

export default {
  getMessages(threadID) {
    if(threadID === 0) {
      return;
    }

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
    if(threadID === 0) {
      return;
    }

    _messageSocket = null;

    var domain;
    let url = window.location.href;

    //find & remove protocol (http, ftp, etc.) and get domain
    if (url.indexOf("://") > -1) {
        domain = url.split('/')[2];
    }
    else {
        domain = url.split('/')[0];
    }

    var sendInterval, reopenInterval;

    var start = function() {
      if(!_messageSocket) {
        _messageSocket = new WebSocket("wss://" + domain + Constants.APIROOT + "/messages/threads/" + threadID + "/socket");
      } else {
        return;
      }

      _messageSocket.onmessage = (res) => {
        let data = JSON.parse(res.data);
        if(data.messages) {
          Dispatcher.handleViewAction({
            type: Constants.ActionTypes.GET_NEW_MESSAGES,
            data: data.messages
          });
        }
      };

      _messageSocket.onopen = () => {
        window.clearInterval(reopenInterval);

        let lastNow = moment(Date.now()).unix();

        sendInterval = window.setInterval(() => {
          if(!_messageSocket || _messageSocket.readyState === 2 || _messageSocket.readyState === 3) {
            window.clearInterval(sendInterval);
            return;
          }

          _messageSocket.send(lastNow.toString());
          lastNow = moment(Date.now()).unix();
        }, 500);
      };

      _messageSocket.onclose = () => {
        reopenInterval = window.setInterval(() => {
          if(!_messageSocket || _messageSocket.readyState === 0 || _messageSocket.readyState === 1) {
            window.clearInterval(reopenInterval);
            return;
          }

          _messageSocket = null;
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
  },

  getMessageCount() {
    var domain;
    let url = window.location.href;

    //find & remove protocol (http, ftp, etc.) and get domain
    if (url.indexOf("://") > -1) {
        domain = url.split('/')[2];
    }
    else {
        domain = url.split('/')[0];
    }

    var _socket = new WebSocket("wss://" + domain + Constants.APIROOT + "/messages/count");

    _socket.onmessage = (res) => {
      let data = JSON.parse(res.data);
      Dispatcher.handleViewAction({
        type: Constants.ActionTypes.GET_MESSAGE_COUNT,
        count: data.count
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
