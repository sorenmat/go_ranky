define(function (require) {

  'use strict';

  /**
   * Module dependencies
   */

  var players  = require('component/players');

  /**
   * Module exports
   */

  return initialize;

  /**
   * Module function
   */

  function initialize() {
      console.log("Flight app starting");
      console.log(players);
    players.attachTo("#flightcontent");
    player.attachTo(document)
  }

});
