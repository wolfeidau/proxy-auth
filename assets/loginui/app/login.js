import React from "react";

export default React.createClass({
  render: function() {
    return (
      <div className="github">
        <a href="/auth/github/authorize">Login With GitHub</a>
      </div>
    );
  },
});
