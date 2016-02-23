import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _questions = [];
let _index = 0;

// Facebook style store creation.
const QuestionsStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  getQuestion() {
    return _questions[_index];
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.GET_QUESTIONS_SUCCESS:
        _questions = action.data;
        QuestionsStore.emitChange();
        break;
      case Constants.ActionTypes.QUESTION_ANSWERED_SUCCESS:
        _index += 1;
        QuestionsStore.emitChange();
        break;
    }
  })
});

export default QuestionsStore;
