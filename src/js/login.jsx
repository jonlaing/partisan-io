import React from 'react';
import jQuery from 'jquery';

global.$ = jQuery;

import Login from './components/Login.jsx';
import Feed from './components/Feed.jsx';

React.render(<Login />, document.getElementById('login'));
React.render(<Feed />, document.getElementById('feed'));
