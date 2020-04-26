import React from 'react';
import './poker.css';
import card_background from './card_background.png'

const Colors = Object.freeze({
  RED: Symbol("red"),
  BLACK: Symbol("black"),
})

const Suits = Object.freeze({
  SPADES: { unicode: "\u2660", color: Colors.BLACK },
  HEARTS: { unicode: "\u2665" , color: Colors.RED },
  DIAMONDS: { unicode: "\u2666" , color: Colors.RED },
  CLUBS: { unicode: "\u2663" , color: Colors.BLACK },
})

function Card(props) {
  if (props.faceDown) {
    return (
      <div class="card-small">
        <img src={card_background} class="card-back" alt="Back of Card"/>
      </div>
    );
  }

  var cardJsx = [];
  if (props.suit.color === Colors.RED) {
    cardJsx.push(<p class="card-value red">{props.value}</p>);
    cardJsx.push(<p class="card-suit red">{props.suit.unicode}</p>);
  } else {
    cardJsx.push(<p class="card-value black">{props.value}</p>);
    cardJsx.push(<p class="card-suit black">{props.suit.unicode}</p>);
  }
  
  return (
    <div class="card-small">
        <div>
          {cardJsx}
        </div>
        <div class="card-bottom">
          {cardJsx}
        </div>
    </div>
  );
}

class Player extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      cards: props.cards,
      chips: props.chips,
      name: props.name,
    };
  }

  dealCard(card) {
    this.setState((state, props) => ({
      cards: state.cards.concat(card)
    }));
  }

  fold() {
    this.setState((state, props) => ({
      cards: []
    }))
  }

  getSeat() {
    return this.state.seat;
  }

  render() {
    return (
      <div class="player-container">
        <div class="player-cards">
          {this.state.cards.map((card, index) =>{
            return (
              <Card
                faceDown={card.faceDown}
                folded={card.folded}
                value={card.value}
                suit={card.suit}
              />
            );
          })}
        </div>
        <div class="player-dashboard">
          <div class="player-name">{this.state.name}</div>
          <div class="player-chips">{this.numberWithCommas(this.state.chips)}</div>
        </div>
      </div>
    )
  }

  numberWithCommas(x) {
    return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  }
}

class ActionCenter extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return <div class="action-center"></div>;
  }
}

class PokerTable extends React.Component {
  constructor(props) {
    super(props);
    this.seatStyles = [
      {
        "position": "absolute",
        "left": "30%",
        "top": "5%",
      },
      {
        "position": "absolute",
        "right": "30%",
        "top": "5%",
      },
      {
        "position": "absolute",
        "right": "3%",
        "top": "25%",
      },
      {
        "position": "absolute",
        "right": "3%",
        "top": "55%",
      },
      {
        "position": "absolute",
        "right": "30%",
        "bottom": "5%",
      },
      {
        "position": "absolute",
        "left": "30%",
        "bottom": "5%",
      },
      {
        "position": "absolute",
        "left": "3%",
        "top": "55%",
      },
      {
        "position": "absolute",
        "left": "3%",
        "top": "25%",
      },
    ];

    this.numSeats = 8;

    const players = {};
    for (var i = 0; i < this.numSeats; ++i) {
      var cards = [];
      if (i % 2 === 1) {
        cards = [{faceDown: true}, {faceDown: true}];
      } else {
        cards = [];
      }

      players[i] = new Player({
        cards: cards,
        chips: 1000,
        name: `Player ${i + 1}`,
      });
    }

    this.state = {
      emptySeats: {
        0: {}
      },
      players: players,
      actionCenter: new ActionCenter({}),
    }
  }

  renderPlayer(seat, player) {
    return (
      <div style={this.seatStyles[seat]}>
        {player.render()}
      </div>
    )
  }

  render() {
      var playerJsx = [];
      for (var seat in this.state.players) {
        playerJsx.push(this.renderPlayer(seat, this.state.players[seat]));
      }

      return (
          <div class="poker-table-container">
            <div class="poker-table"></div>
            {playerJsx}
            {this.state.actionCenter.render()}
          </div>
      );
  }
}

export default PokerTable;