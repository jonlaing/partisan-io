import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  getThreads() {
    $.ajax({
      url: Constants.APIROOT + '/messages/threads',
      method: 'GET',
      dataType: 'json'
    })
      .done(function(res) {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.GET_THREADS_SUCCESS,
          data: res
        });
      })
      .fail(function(res) {
        console.log(res);
      });
  },

  createThread(friendID) {
    $.ajax({
      url: Constants.APIROOT + '/messages/threads',
      method: 'POST',
      data: { user_id: friendID },
      dataType: 'json'
    })
    .done((res) => {
      Dispatcher.handleViewAction({
        type: Constants.ActionTypes.CREATE_THREAD_SUCCESS,
        thread: res.thread
      });
    })
    .fail(function(res) {
      console.log(res);
    });
  },

  switchThreads(threadID) {
    Dispatcher.handleViewAction({
      type: Constants.ActionTypes.SWITCH_THREADS,
      threadID: threadID
    });
  }

};
