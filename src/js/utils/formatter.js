import marked from 'marked';


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
    let content = marked(body, {renderer: _markedRenderer} );
    content = _hashtagify(content);
    content = _tagify(content);
    return { __html: content };
  }
};
