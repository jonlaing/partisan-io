/*global data */
import React from 'react';
import ReactDOM from 'react-dom';
import jQuery from 'jquery';

import Events from 'events';
Events.EventEmitter.prototype._maxListeners = 100;

global.$ = jQuery;

import FrontTicker from './components/FrontTicker.jsx';
import PasswordReset from './components/PasswordReset.jsx';

// optionally attach DOM elements to React
let frontTicker = document.getElementById('front-ticker');
let passwordReset = document.getElementById('password-reset');

if(frontTicker != null) {
  ReactDOM.render(<FrontTicker />, frontTicker);
}

if(passwordReset != null) {
  ReactDOM.render(<PasswordReset resetID={data.reset_id}/>, passwordReset);
}
