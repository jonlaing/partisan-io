/*global $, WebSocket */
import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

import moment from 'moment';

export default {
  getFeed() {
    $.ajax({
      url: '/api/v1/feed',
      dataType: 'json',
      method: 'GET'
    })
      .done(function(res) {
        let data = res.feed_items;

        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_FEED,
          data: data
        });
      })
      .fail(function(res) {
        // Logged out
        if(res.status === 401) {
          Dispatcher.handleViewAction({
            type: Constants.ActionTypes.LOGGED_OUT
          });
        } else if (res.status === 404) {
          Dispatcher.handleViewAction({
            type: Constants.ActionTypes.NO_FRIENDS
          });
        }
      });
  },
  getPage(page) {
    var p = page || 1;

    $.ajax({
      url: '/api/v1/feed',
      data: { page: p },
      dataType: 'json',
      method: 'GET'
    })
      .done(function(res) {
        let data = res.feed_items;

        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_FEED_PAGE,
          data: data
        });
      })
      .fail(function(res) {
        // Logged out
        if(res.status === 401) {
          Dispatcher.handleViewAction({
            type: Constants.ActionTypes.LOGGED_OUT
          });
        }
      });
  },
  getByUser(userID) {
    $.ajax({
      url: '/api/v1/feed/show/' + userID,
      dataType: 'json',
      method: 'GET'
    })
      .done(function(res) {
        let data = res.feed_items;

        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_FEED,
          data: data
        });
      })
      .fail(function(res) {
        // Logged out
        if(res.status === 401) {
          Dispatcher.handleViewAction({
            type: Constants.ActionTypes.LOGGED_OUT
          });
        }
      });
  },
  feedSocket() {
    var domain;
    let url = window.location.href;

    //find & remove protocol (http, ftp, etc.) and get domain
    if (url.indexOf("://") > -1) {
        domain = url.split('/')[2];
    }
    else {
        domain = url.split('/')[0];
    }

    var _socket = new WebSocket("ws://" + domain + Constants.APIROOT + "/feed/socket");

    _socket.onmessage = (res) => {
      let data = JSON.parse(res.data);
      Dispatcher.handleViewAction({
        type: Constants.ActionTypes.GET_NEW_FEED_ITEMS,
        data: data.feed_items
      });
    };

    _socket.onopen = () => {
      let lastNow = moment(Date.now()).unix();
      window.setInterval(() => {
        _socket.send(lastNow.toString());
        lastNow = moment(Date.now()).unix();
      }, 600000);
    };
  },
  addItem(data) {
    Dispatcher.handleViewAction({
      type: Constants.ActionTypes.ADD_FEED_ITEM,
      data: data
    });
  }
};
