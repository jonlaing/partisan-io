import React from 'react';

import moment from 'moment';

import MatchesActionCreator from '../actions/MatchesActionCreator';

import Card from './Card.jsx';

export default React.createClass({
  getInitialState() {
    return {
      miles: 25,
      gender: "",
      minAge: 18,
      maxAge: moment().diff(this.props.birthdate, 'years') + 10,
      changed: false
    };
  },

  onMilesChange(e) {
    this.setState({miles: e.target.value, changed: true});
  },

  onRedoSearch() {
    let minAge = this._getAge(this.state.minAge);
    let maxAge = this._getAge(this.state.maxAge);
    if(minAge > this.state.maxAge) {
      minAge = this.state.maxAge;
    }
    if(maxAge < this.state.minAge) {
      maxAge = this.state.minAge;
    }

    this.setState({changed: false});
    MatchesActionCreator.getMatches(this.state.miles, this.state.gender, minAge, maxAge);
  },

  onMinAgeChange(e) {
    this.setState({minAge: e.target.value, changed: true});
  },

  onMinAgeBlur(e) {
    let age = this._getAge(e.target.value);
    if(age > this.state.maxAge) {
      age = this.state.maxAge;
    }

    this.setState({minAge: age});
  },

  onMaxAgeChange(e) {
    this.setState({maxAge: e.target.value, changed: true});
  },

  onMaxAgeBlur(e) {
    let age = this._getAge(e.target.value);
    if(age < this.state.minAge) {
      age = this.state.minAge;
    }

    this.setState({maxAge: age});
  },

  onGenderChange(e) {
    this.setState({gender: e.target.value, changed: true});
  },

  componentDidMount() {
  },

  render() {
    var button;

    if(this.state.changed === true) {
      button = <button className="button button-small" onClick={this.onRedoSearch}>Redo Search</button>;
    }

    return (
      <div className="matches-search">
        <Card>
          <h2>Search Parameters</h2>
          <div className="matches-search-distance">
            <label>Distance:</label>
            <input id="miles" name="miles" type="range" min={0} max={100} defaultValue={25} onChange={this.onMilesChange}/>
            <output htmlFor="miles">{this.state.miles} miles</output>
          </div>
          <div className="matches-search-ages">
            <div>
              <label>Min Age:</label>
              <input type="number" value={this.state.minAge} onChange={this.onMinAgeChange} onBlur={this.onMinAgeBlur} />
            </div>
            <div>
              <label>Max Age:</label>
              <input type="number" value={this.state.maxAge} onChange={this.onMaxAgeChange} onBlur={this.onMaxAgeBlur} />
            </div>
          </div>
          <div className="matches-search-gender">
            <label>Gender:</label>
            <input type="text" placeholder="Any gender" onChange={this.onGenderChange} />
          </div>
          <div className="text-center">
            {button}
          </div>
        </Card>
      </div>
    );
  },

  _getAge(val) {
    let age = parseInt(val);
    // can't search for minors!
    if(age < 18) {
      age = 18;
    }

    return age;
  }
});
