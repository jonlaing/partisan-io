import React from 'react';

let _ENTER = 13; // key code for pressing the ENTER/RETURN key

export default React.createClass({
  getInitialState() {
    return {
      editLocation: false,
      editGender: false
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

  componentDidMount() {
  },

  render() {
    var location, gender;

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

    return (
      <div className="profile-edit-info">
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
        <input type="text" defaultValue={this.props.postalCode} onKeyDown={this.handleLocationKeyDown} />
      </div>
    );
  },
  _editGenderTemplate() {
    let g = this.props.gender || "";
    return (
      <div>
        <input type="text" placeholder="Type in your gender" defaultValue={g} onKeyDown={this.handleGenderKeyDown} />
      </div>
    );
  },

  _cityState(location) {
    return location.replace(/\s\d+.*$/, '');
  }
});
