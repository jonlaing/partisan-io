import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _questions = [];

// add private functions to modify data
function addQuestion(question) {
  _questions.push(question);
}

// Facebook style store creation.
const QuestionsStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  getLast() {
    return _questions[_questions.length - 1];
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.GET_QUESTION_SUCCESS:
        let question = action.data;
        // TODO: check for unique
        addQuestion(question);
        QuestionsStore.emitChange();
        break;
    }
  })
});

export default QuestionsStore;
