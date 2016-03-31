import Dispatcher from '../Dispatcher';
import Constants from '../Constants';

export default {
  beginReport(id, type) {
    Dispatcher.handleViewAction({
      type: Constants.ActionTypes.BEGIN_FLAG,
      id: id,
      recordType: type
    });
  },

  cancelReport() {
    Dispatcher.handleViewAction({
      type: Constants.ActionTypes.CANCEL_FLAG
    });
  },

  submitReport(id, type, reason, message) {
    console.log(id, type, reason, message);
    $.ajax({
      url: Constants.APIROOT + '/flag',
      data: JSON.stringify({
        record_id: id,
        record_type: type,
        reason: reason,
        message: message
      }),
      method: 'POST',
      dataType: 'json'
    })
      .fail(function(res) {
        console.log(res);
      })
      .always(function() {
        Dispatcher.handleViewAction({
          type: Constants.ActionTypes.SUBMIT_FLAG
        });
      });
  }
};
