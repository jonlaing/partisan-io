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
  }
};
