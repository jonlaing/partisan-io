import React from 'react';

import LightboxActionCreator from '../actions/LightboxActionCreator';
import LightboxStore from '../stores/LightboxStore';

import Modal from './Modal.jsx';

export default React.createClass({
  getInitialState() {
    return { show: false, image: "https://s3-us-west-2.amazonaws.com/partisan-staging/img/QFBhJyMDyHAdzhfP.jpg" };
  },

  handleClose() {
    LightboxActionCreator.close();
  },

  componentDidMount() {
    LightboxStore.addChangeListener(this.onChange);
  },

  componentWillUnmount() {
    LightboxStore.removeChangeListener(this.onChange);
  },

  render() {
    return (
      <Modal show={this.state.show} className="lightbox" onCloseClick={this.handleClose}>
        <img src={this.state.image} className="lightbox-image" />
      </Modal>
    );
  },

  onChange() {
    this.setState(LightboxStore.getState);
  }
});
