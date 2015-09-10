import React from 'react';
import QuestionsActionCreator from '../actions/QuestionsActionCreator';
import QuestionsStore from '../stores/QuestionsStore';
import Card from './Card.jsx';

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
        <div className="row">
          <div className="large-12 columns">
            <div className="right">
              {this.state.questionsAnswered} of {this._maxQuestions()}
            </div>
          </div>
        </div>
        <div className="row">
          <div className="large-12 columns">
            <Card>
              <div className="card-body">
                {this.state.question.prompt}
              </div>
            </Card>
          </div>
        </div>
        <div className="row">
          <div className="large-4 columns">
            <button className="button alert expand" onClick={this.handleDisagree}>Disagree</button>
          </div>
          <div className="large-4 large-offset-4 columns">
            <button className="button success expand" onClick={this.handleAgree}>Agree</button>
          </div>
        </div>
      </div>
    );
  },

  _onChange() {
    let question = QuestionsStore.getLast();
    let answered = this.state.questionsAnswered + 1;


    if(this._maxQuestions() > -1 && answered > this._maxQuestions()) {
      window.location.href = "/feed.html";
      return;
    }

    console.log(question);
    this.setState({question: question, questionsAnswered: answered});
  },

  _maxQuestions() {
    if(this.props.maxQuestions === undefined) {
      return 15;
    }

    return this.props.maxQuestions;
  }
});
