import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

import QuestionsActionCreator from '../actions/QuestionsActionCreator';

// data storage
let _questions = [];
let _index = 0;

const maxQuestions = 20;

// Facebook style store creation.
const QuestionsStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  getQuestion() {
    console.log("questions:", _questions);
    return _questions[_index];
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    switch(action.type) {
      case Constants.ActionTypes.GET_QUESTIONS_SUCCESS:
        _questions = _questions.concat(action.data.questions);
        QuestionsStore.emitChange();
        break;
      case Constants.ActionTypes.QUESTION_ANSWERED_SUCCESS:
        _index += 1;

        // if we're answered a multiple of 4, then we need new questions
        // otherwise emit the change
        if(_index % 4 === 0 && _index !== maxQuestions) {
          console.log("getting new questions");
          QuestionsActionCreator.getQuestions();
        } else {
          console.log("index:", _index);
          QuestionsStore.emitChange();
        }
        break;
    }
  })
});

export default QuestionsStore;
