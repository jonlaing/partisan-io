import React from 'react';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';

var _timer;

const _fullList = [
  "liberals",
  "conservatives",
  "anarchists",
  "communists",
  "libertarians",
  "Democrats",
  "Republicans",
  "socialists",
  "progressives",
  "marxists",
  "activists",
  "conspiracy theorists",
  "feminists",
  "traditionalists",
  "radicals",
  "patriots",
  "everyone!"
];

export default React.createClass({
  getInitialState() {
    return {
      current: [_fullList[0]]
    };
  },

  componentDidMount() {
    _timer = window.setInterval(this._cycle, 2000);
  },

  componentWillUnmount() {
    window.clearInterval(_timer);
  },

  onWindowBlur() {
    window.clearInterval(_timer);
  },

  onWindowFocus() {
    _timer = window.setInterval(this._cycle, 2000);
  },

  render() {
    let text = this.state.current.map((t, i) => {
      if(i === 0) { return <h3 key={Date.now()}>{t}</h3>; }
    });

    return (
      <div className="front-ticker">
        <h3>Partisan is for:</h3>
        <div className="front-ticker-tick">
          <ReactCSSTransitionGroup transitionName="front-ticker-text" transitionEnterTimeout={2000} transitionLeaveTimeout={2000}>
            {text}
          </ReactCSSTransitionGroup>
        </div>
      </div>
    );
  },

  _cycle() {
      let i = Math.floor(Math.random() * _fullList.length);
      if(_fullList[i] === this.state.current[0]) {
        return this._cycle();
      }

      this.setState({current: [_fullList[i]]});
  }
});
