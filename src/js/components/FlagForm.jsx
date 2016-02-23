import React from 'react';
import RadioGroup from 'react-radio';

import FlagActionCreator from '../actions/FlagActionCreator';

import Modal from './Modal.jsx';

let _reasons = [
  { value: 1, label: 'Offensive Content'},
  { value: 2, label: 'Copyrighted Material'},
  { value: 3, label: 'Spam' },
  { value: 4, label: 'Other' }
];

export default React.createClass({
  getInitialState() {
    return { reason: 0 };
  },

  handleClose() {
    FlagActionCreator.cancelReport();
  },

  handleReasonChange(val) {
    this.setState({ reason: val });
  },

  handleSubmit() {
    let message = $(React.findDOMNode(this.refs.message)).val();
    FlagActionCreator.submitReport(this.props.id, this.props.type, parseInt(this.state.reason), message);
  },

  render() {
    if(this.props.show === true) {
      return (
        <Modal show={this.props.show} onCloseClick={this.handleClose}>
          <h4>Flag a {this.props.type}</h4>
          <hr/>
          <RadioGroup name="reason" items={_reasons} onChange={this.handleReasonChange} />
          <div>
            <label>Comment:</label>
            <textarea name="message" ref="message" placeholder="Type any additional comments..."/>
          </div>
          <div>
            <button className="button" onClick={this.handleSubmit}>Submit</button>
          </div>
        </Modal>
      );
    } else {
      return (
        <Modal show={this.props.show} onCloseClick={this.handleClose}>
        </Modal>
      );
    }
  }
});
