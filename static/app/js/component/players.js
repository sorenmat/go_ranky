define(function (require) {

  'use strict';

  /**
   * Module dependencies
   */

  var defineComponent = require('flight/lib/component'),
    withAjax = require('flight-tower/lib/with_ajax');

  /**
   * Module exports
   */

  return defineComponent(players, withAjax);

  /**
   * Module function
   */

  function players() {
/*    this.attributes({

    });
*/
    this.handlePlayers = function(event, data) {
	var players = data.players;
  var htmlCode = ""
	for (var i = 0; i < players.length; i++) {
	htmlCode = htmlCode + " <div class='white-panel pn'> \
 <div class='white-header'> \
 <h5>#"+(i+1)+"</h5> \
             </div> \
             <p><img src='img/ui-zac.jpg' class='img-circle' width='50'></p> \
             <p><b>"+players[i].Name+"</b></p> \
               <div class='row'> \
                 <div class='col-md-6'> \
                   <p class='small mt'>MEMBER SINCE</p> \
                   <p>2012</p> \
                 </div> \
                 <div class='col-md-6'> \
                   <p class='small mt'>TOTAL SPEND</p> \
                   <p>$ 47,60</p> \
                 </div> \
               </div> \
           </div> "
	}
this.$node.html(htmlCode)
	console.log(data.players[1]);

    }
    this.handlePlayersFail = function(event, data) {
	console.log("Dohhh");

    }
      this.getPlayers = function() {
	  var that = this;
	  $.getJSON( "/players", function( data ) {
	      console.log(data);
	      that.trigger('dataPlayers', {players: data})
	  });
      }

    this.after('initialize', function () {
	this.on('requestPlayers', this.getPlayers);
	this.on(document, 'dataPlayers', this.handlePlayers);
	this.on(document, 'dataPlayersError', this.handlePlayersFail);
	this.get('/players', 'Players');
	this.trigger('requestPlayers');
    });
  }

});
