import React from 'react';
import CheckboxGroup from 'react-checkbox-group';

export default React.createClass({
  handleChange() {
    var bitMap = 0;
    let values = this.refs.lookingFor.getCheckedValues();

    values.forEach(function(v) {
      bitMap += parseInt(v);
    });

    this.props.onChange(bitMap);
  },

  componentDidMount() {
  },

  render() {
    let lookingFor = this._parseLookingFor(this.props.lookingFor);

    return (
      <div className="lookingfor">
        <CheckboxGroup
          name="looking_for"
          value={lookingFor}
          ref="lookingFor"
          onChange={this.handleChange} >
            <label className={(this._active(1 << 0) ? " active" : "")}>
              <i className="fi-torsos-all-female"></i>
              <input type="checkbox" value={1 << 0} /> Friends
            </label>
            <label className={(this._active(1 << 1) ? " active" : "")}>
              <i className="fi-heart"></i>
              <input type="checkbox" value={1 << 1} /> Love
            </label>
            <label className={(this._active(1 << 2) ? " active" : "")}>
              <i className="fi-skull"></i>
              <input type="checkbox" value={1 << 2} /> Enemies
            </label>
        </CheckboxGroup>
      </div>
    );
  },

  _parseLookingFor(n) {
    var vals = [];
    for(var i = 0; i <= 3; i++) {
      if((n & 1 << i) !== 0) {
        vals.push((1 << i).toString());
      }
    }
    return vals;
  },

  _active(n) {
    let vals = this._parseLookingFor(this.props.lookingFor);
    for(let i in vals) {
      if(parseInt(vals[i]) === n) {
        return true;
      }
    }
    return false;
  }

});
