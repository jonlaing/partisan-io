import React from 'react';
import QuestionsActionCreator from '../actions/QuestionsActionCreator';
import QuestionsStore from '../stores/QuestionsStore';
import Card from './Card.jsx';
import UserSession from './UserSession.jsx';

export default React.createClass({
  getInitialState() {
    return { question: { prompt: "Fetching...", map: [] }, questionsAnswered: 0};
  },

  handleAgree() {
    QuestionsActionCreator.answerQuestion(this.state.question, true);
  },

  handleDisagree() {
    QuestionsActionCreator.answerQuestion(this.state.question, false);
  },

  componentDidMount() {
    QuestionsStore.addChangeListener(this._onChange);
    QuestionsActionCreator.getQuestion();
  },

  componentWillUnmount() {
    QuestionsStore.removeChangeListener(this._onChange);
  },

  render() {
    return (
      <div className="question">
        <header>
          <UserSession username={this.props.data.user.username} />
        </header>
        <div className="question-container">
          <div className="question-number">
            {this.state.questionsAnswered} of {this._maxQuestions()}
          </div>
          <div className="question-body">
            <Card>
              <div className="card-body">
                {this.state.question.prompt}
              </div>
            </Card>
          </div>
        </div>
        <div className="question-actions">
          <button className="button alert expand" onClick={this.handleDisagree}>Disagree</button>
          <button className="button success expand" onClick={this.handleAgree}>Agree</button>
        </div>
      </div>
    );
  },

  _onChange() {
    let question = QuestionsStore.getLast();
    let answered = this.state.questionsAnswered + 1;


    if(this._maxQuestions() > -1 && answered > this._maxQuestions()) {
      window.location.href = "/profile";
      return;
    }

    this.setState({question: question, questionsAnswered: answered});
  },

  _maxQuestions() {
    if(this.props.maxQuestions === undefined) {
      return 15;
    }

    return this.props.maxQuestions;
  }
});
