import React from 'react';
import moment from 'moment';

let _ENTER = 13; // key code for pressing the ENTER/RETURN key

export default React.createClass({
  getInitialState() {
    return {
      editLocation: false,
      editGender: false,
      editBirthdate: false
    };
  },

  handleLocationClick() {
    this.setState({editLocation: true});
  },

  handleLocationKeyDown(e) {
    if(e.keyCode === _ENTER && e.target.value !== "") {
      this.setState({editLocation: false});
      this.props.onLocationFinish(e);
    }
  },

  handleGenderClick() {
    this.setState({editGender: true});
  },

  handleGenderKeyDown(e) {
    if(e.keyCode === _ENTER && e.target.value !== "") {
      this.setState({editGender: false});
      this.props.onGenderFinish(e);
    }
  },

  handleBirthdateClick() {
    this.setState({editBirthdate: true});
  },

  handleBirthdateChange(e) {
    console.log(e.target.value);
    this.setState({editBirthdate: false});
    this.props.onBirthdateFinish(e);
  },

  componentDidMount() {
  },

  render() {
    var location, gender, birthdate;

    if(this.state.editLocation === false) {
      var cityState = this._cityState(this.props.location);

      location = (
        <div>
          <a href="javascript:void(0)" onClick={this.handleLocationClick}>{cityState}</a>
        </div>
      );
    } else {
      location = this._editLocationTemplate();
    }

    if(this.state.editGender === false) {
      let g = this.props.gender || "No Gender";
      gender = (
        <div>
          <a href="javascript:void(0)" onClick={this.handleGenderClick}>{g}</a>
        </div>
      );
    } else {
      gender = this._editGenderTemplate();
    }

    if(this.state.editBirthdate === false) {
      var b;
      if(moment(this.props.birthdate).isBefore('1800-12-31')) {
        b = "No Age";
      } else {
        b = moment().diff(this.props.birthdate, 'years') + " years old";
      }

      birthdate = (
        <div>
          <a href="javascript:void(0)" onClick={this.handleBirthdateClick}>{b}</a>
        </div>
      );
    } else {
      birthdate = this._editBirthdateTemplate();
    }

    return (
      <div className="profile-edit-info">
        <div className="profile-edit-birthdate">
          {birthdate}
        </div>
        <div className="profile-edit-location">
          {location}
        </div>
        <div className="profile-edit-gender">
          {gender}
        </div>
      </div>
    );
  },

  _editLocationTemplate() {
    return (
      <div>
        <input type="text" defaultValue={this.props.postalCode} onKeyDown={this.handleLocationKeyDown} autoFocus={true} />
      </div>
    );
  },
  _editGenderTemplate() {
    let g = this.props.gender || "";
    return (
      <div>
        <input type="text" placeholder="Type in your gender" defaultValue={g} onKeyDown={this.handleGenderKeyDown} autoFocus={true} />
      </div>
    );
  },
  _editBirthdateTemplate() {
    var date = moment(this.props.birthdate);

    if(parseInt(date.get('year')) === 0) {
      date = "";
    }

    return (
      <div>
        <input type="date" defaultValue={date} onChange={this.handleBirthdateChange} />
      </div>
    );
  },


  _cityState(location) {
    return location.replace(/\s\d+.*$/, '');
  }
});
