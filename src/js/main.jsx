import React from 'react';
import jQuery from 'jquery';

global.$ = jQuery;

import Login from './components/Login.jsx';
import SignUp from './components/SignUp.jsx';
import UserSession from './components/UserSession.jsx';
import Feed from './components/Feed.jsx';

// optionally attack DOM elements to React
let login = document.getElementById('login');
let signUp = document.getElementById('sign-up');
let userSession = document.getElementById('user-session');
let feed = document.getElementById('feed');

// for static login page
if(login !== null) {
  React.render(<Login />, login);
}

if(signUp !== null) {
  React.render(<SignUp />, signUp);
}

if(userSession !== null) {
  React.render(<UserSession />, userSession);
}

if(feed !== null) {
  React.render(<Feed />, feed);
}
