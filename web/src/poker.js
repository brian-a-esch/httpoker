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
    )
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

class Seat extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      cards: [],
      chips: props.chips,
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

  render() {
    return (
      <div class="seat">
        <div>
          {this.state.cards.map((card, index) =>{
            return <Card card/>
          })}
        </div>
        <div class="chip-count"></div>
      </div>
    )
  }
}

class PokerTable extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
      return (
          <div>
            <div>
                <Card
                  faceDown={true}
                  value="K"
                  suit={Suits.CLUBS}
                />
                <Card
                  faceDown={false}
                  value="K"
                  suit={Suits.CLUBS}
                />
            </div>
            <div class="poker-table"></div>
          </div>
      );
  }
}

export default PokerTable;