import React from 'react';
import marked from 'marked';
import moment from 'moment';


marked.setOptions({
  sanitize: true,
  tables: false
});

var _markedRenderer = new marked.Renderer();

_markedRenderer.heading = (text) => {
  return '<p><strong>' + text + '</strong></p>';
};

function _hashtagify(content) {
  let tags = content.match(/#[a-zA-Z]+/g);
  var newContent = content;

  if(tags !== null) {
    tags.forEach((tag) => {
      let encoded = encodeURIComponent(tag);
      newContent = newContent.replace(tag, '<a href="/hashtags?q=' + encoded + '">' + tag + '</a>');
    });
  }

  return newContent;
}

function _tagify(content) {
  let tags = content.match(/@([a-zA-Z0-9_.]+)/g); // will match emails too, will filter later
  var newContent = content;

  if(tags !== null) {
    tags.forEach((tag) => {
      // make sure this isn't an email
      if(tag.match(/\./) === null) {
        newContent = newContent.replace(tag, '<a href="/profiles/' + tag.substring(1) + '">' + tag + '</a>');
      }
    });
  }

  return newContent;
}

export default {
  // format for a Post
  post(body) {
    let content = marked(body, {renderer: _markedRenderer} );
    content = _hashtagify(content);
    content = _tagify(content);
    return { __html: content };
  },

  // format for a Comment
  comment(body) {
    let content = marked(body, {renderer: _markedRenderer} );
    content = _hashtagify(content);
    content = _tagify(content);
    return { __html: content };
  },

  // format for a User Summary
  userSummary(body) {
    if(body === undefined || body.length < 1) {
      return { __html: '<em>No Summary</em>' };
    }

    let content = marked(body, {renderer: _markedRenderer} );
    content = _hashtagify(content);
    content = _tagify(content);
    return { __html: '<div>' + content + '</div>' };
  },

  cityState(location) {
    return location.replace(/\s\d+.*$/, '');
  },

  age(birthdate, showIfNone) {
    // You're probably not older than 200 years old, which means
    // your bithdate hasn't been entered
    if(moment(birthdate).isBefore('1800-12-31')) {
      if(showIfNone === false) {
        return "";
      }

      return "No Age";
    }

    // return moment().diff(birthdate, 'years') + " years old";
    return moment().diff(birthdate, 'years') + "yrs";
  },

  avatarUrl(url) {
    if(url.length < 1) {
      return "/images/avatar.jpg";
    }

    return url;
  },

  match(percent) {
    var height;

    if(percent >= 85) {
      height = "high";
    } else if(percent < 85 && percent >= 65) {
      height = "med";
    } else {
      height = "low";
    }

    return <span className={height}>{percent}% Match</span>;
  }
};
