import React from 'react';
import jQuery from 'jquery';

global.$ = jQuery;

import Login from './components/Login.jsx';
import Feed from './components/Feed.jsx';

// optionally attack DOM elements to React
let login = document.getElementById('login');
let feed = document.getElementById('feed');

// for static login page
if(login !== null) {
  React.render(<Login />, login);
}

if(feed !== null) {
  React.render(<Feed />, feed);
}
