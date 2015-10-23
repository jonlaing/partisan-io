import React from 'react';

const _months = [
  [ 'Jan', 31],
  [ 'Feb', 29],
  [ 'March', 31],
  [ 'April', 30],
  [ 'May', 31],
  [ 'June', 30],
  [ 'July', 31],
  [ 'Aug', 31],
  [ 'Sept', 30],
  [ 'Oct', 31],
  [ 'Nov', 30],
  [ 'Dec', 31]
];

const _youngestYear = new Date().getFullYear() - 18;

export default React.createClass({
  getInitialState() {
    let matches = this.props.defaultDate.match(/^(\d{4})-(\d{2})-(\d{2})/);
    if (matches !== null) {
      console.log(matches);
      return { selYear: parseInt(matches[1]), selMonth: parseInt(matches[2]), selDay: parseInt(matches[3]), done: false };
    }

    return {selMonth: 1, selDay: 1, selYear: _youngestYear, done: false };
  },

  handleMonthChange(e) {
    this.setState({selMonth: e.target.value});
  },

  handleDayChange(e) {
    this.setState({selDay: e.target.value});
  },

  handleYearChange(e) {
    this.setState({selYear: e.target.value, done: true});
  },

  getDate() {
    let month = (this.state.selMonth < 10) ? "0" + this.state.selMonth : this.state.selMonth + "";
    let day = (this.state.selDay < 10) ? "0" + this.state.selDay : this.state.selDay + "";

    console.log(this.state.selYear);

    return this.state.selYear + "-" + month + "-" + day;
  },

  componentDidUpdate() {
    if(this.state.done === true) {
      this.props.onChange();
    }
  },

  render() {
    var months = _months.map((v, k) => {
      let month = k + 1;
      return <option key={month} value={month}>{v[0]}</option>;
    });

    var days = [];
    for(let i = 1; i <= _months[this.state.selMonth - 1][1]; i++) {
      days.push(<option key={i} value={i}>{i}</option>);
    }

    var years = [];
    for(let j = _youngestYear; j > _youngestYear - 200; j--) {
      years.push(<option key={j} value={j}>{j}</option>);
    }

    return (
      <div className="datepicker">
        <select onChange={this.handleMonthChange} value={this.state.selMonth} ref="month">
          {months}
        </select>
        <select onChange={this.handleDayChange} value={this.state.selDay} ref="day">
          {days}
        </select>
        <select onChange={this.handleYearChange} value={this.state.selYear} ref="year">
          {years}
        </select>
      </div>
    );
  }
});
