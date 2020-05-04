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

function PlayingCard(props) {
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

async function makeRequest(url, method, data) {
  return await fetch(url, {
    method: method, 
    cache: 'no-cache', 
    headers: {
      'Content-Type': 'application/json'
    },
    referrerPolicy: 'no-referrer', 
    body: JSON.stringify(data) 
  });
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
              <PlayingCard
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
    this.state = {};
  }

  render() {
    return <div class="action-center"></div>;
  }
}

async function getGameState(gameID, passphrase) {
  let r = await makeRequest("/api/v1/game/status", "POST", {gameID: gameID, passphrase: passphrase});
  if (r.ok) {
    let game = await r.json();
    return game;
  } 
  let error = await r.text();
  return {error: error};
}

class GameLogin extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      passphrase: '',
      error: '',
    };
  }

  handlePassphraseChange = (event) => {
    this.setState({passphrase: event.target.value});
  }

  handleSubmit = async (event) => {
    event.preventDefault();
    const passphrase = this.state.passphrase;
    let gameStateMaybe = await getGameState(this.props.gameID, passphrase)
    if ('error' in gameStateMaybe) {
      this.setState({error: gameStateMaybe.error});
      return;
    }

    this.props.handleGameLogin(gameStateMaybe, passphrase);
  }

  render() {
    return (
      <div style={{width: '400px', position: 'absolute'}} class="center">
        <form onSubmit={this.handleSubmit}>
          <label>
            Passphrase:
            <input type="text" value={this.state.passphrase} onChange={this.handlePassphraseChange}></input>
          </label>
          <input type="submit" value="Login"/>
          <text style={{color: 'red'}} >{this.state.error}</text>
        </form>
      </div>
    );
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

    this.state = {
      emptySeats: [],
      players: {},
      passphrase: props.initialPassphrase,
      blindSize: 0,
      error: '',
    }
  }

  renderPlayer(seat, player) {
    return (
      <div style={this.seatStyles[seat]}>
        {player.render()}
      </div>
    )
  }

  handleGameLogin = (gameState, passphrase) => {
    this.setState({
      emptySeats: gameState.emptySeats,
      players: gameState.players,
      blindSize: gameState.blindSize,
      passphrase: passphrase,
    })
  }

  async componentDidMount() {
    if (this.loggedIn()) {
      const passphrase = this.state.passphrase;
      let gameStateMaybe = await getGameState(this.state.gameID, passphrase);
      if ('error' in gameStateMaybe) {
        throw Error("Passphrase from game creation should result in login");
      }

      this.handleGameLogin(gameStateMaybe, passphrase)
    }
  }

  loggedIn = () => {
    return this.state.passphrase.length !== 0;
  }

  render() {
    if (this.state.error.length !== 0) {
      return (
      <text style={{color: 'red'}} class="error">An error has occurred: {this.state.error}</text>
      )
    }

    if (!this.loggedIn()) {
      return <GameLogin 
        handleGameLogin={this.handleGameLogin}
        error={this.state.error}
        />
    }

    var playerJsx = [];
    for (var seat in this.state.players) {
      playerJsx.push(this.renderPlayer(seat, this.state.players[seat]));
    }

    return (
      <div class="root-background">
        <h2>{this.state.gameTitle}</h2>
        <div class="poker-table-container">
          <div class="poker-table"></div>
          {playerJsx}
          <ActionCenter/>
        </div>
      </div>
    );
  }
}

class CreateGameComponent extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      passphrase: '',
      starterChips: '',
      blindSize: '',
      error: '',
    }
  }

  handleSubmit = (event) => {
    event.preventDefault();
    const passphrase = this.state.passphrase;
    const starterChips = this.parsePositiveNumber(this.state.starterChips);
    if (starterChips === -1) {
      this.setState({error: 'Starter Chips must be a positive integer'});
      return;
    } 

    const blindSize = this.parsePositiveNumber(this.state.blindSize);
    if (blindSize === -1) {
      this.setState({error: 'Blind Size must be a positive integer'});
      return;
    }

    makeRequest('api/v1/game/create', 'POST', {
      passphrase: passphrase,
      starterChips: starterChips,
      blindSize: blindSize,
    }).then(response => {
      if (response.ok) {
        response.json().then(game => {this.props.onGameCreateSuccess(game, passphrase)})
      } else {
        response.text().then(err => this.setState({error: err}))
      }
    }).catch(error => {
      this.setState({error: 'An error has occurred'});
    })
  }

  parsePositiveNumber = (numberStr) => {
    var n = Math.floor(Number(numberStr));
    if (n !== Infinity && String(n) === numberStr && n >= 0) {
      return n;
    }
    return -1;
  }

  handlePassphraseChange = (event) => {
    this.setState({passphrase: event.target.value});
  }

  handleStarterChipsChange = (event) => {
    this.setState({starterChips: event.target.value});
  }

  handleBindSizeChange = (event) => {
    this.setState({blindSize: event.target.value});
  }

  render() {
    return (
      <div style={{width: '400px', position: 'absolute'}} class="center">
        <form onSubmit={this.handleSubmit}>
          <label>
            Passphrase:
            <input type="text" value={this.state.passphrase} onChange={this.handlePassphraseChange}></input>
          </label>
          <label>
            Starter Chips:
            <input type="text" value={this.state.starterChips} onChange={this.handleStarterChipsChange}></input>
          </label>
          <label>
            Blind Size:
            <input type="text" value={this.state.blindSize} onChange={this.handleBindSizeChange}></input>
          </label>
          <input type="submit" value="Create"/>
          <text style={{color: 'red'}} >{this.state.error}</text>
        </form>
      </div>
    );
  }
}

const Pages = Object.freeze({
  MAIN: "main",
  CREATE_GAME: "create_game",
  GAME: "game",
})

class MainComponent extends React.Component {
  constructor(props) {
    super(props);
    const pathname = window.location.pathname;
    const query = window.location.search;

    var page = null;
    var url = null;
    var gameID = NaN;

    if (pathname === "/create") {
      page = Pages.CREATE_GAME
      url = '/create'
    } else if (pathname === "/game") {
      const params = new URLSearchParams(query);
      gameID = parseInt(params.get('gameID'));
      if (!isNaN(gameID)) {
        page = Pages.GAME;
        url = `/game?gameID=${gameID}`;
      } else {
        page = Pages.MAIN;
        url = ''
      }
    } else {
        page = Pages.MAIN;
        url = ''
    }

    window.history.replaceState({ page: page }, '', url);
    this.state = {
      page: page,
      // Only valid if page is GAME
      gameID: gameID,
      passphrase: '',
    }
  }

  handleCreateClick = () =>{
    window.history.pushState({ page: Pages.CREATE_GAME}, '', '/create');
    this.setState({ page: Pages.CREATE_GAME });
  }

  handlePopState = (event) => {
    this.setState({ page: event.state.page });
  }

  componentDidMount() {
    window.addEventListener("popstate", this.handlePopState);
  }

  componentWillUnmount() {
    window.removeEventListener("popstate", this.handlePopState);
  }

  handleCreateGameSuccess = (game, passphrase) => {
    window.history.pushState({page: Pages.GAME}, '', `/game?gameID=${game.gameID}`)
    this.setState({ page: Pages.GAME, gameID: game.gameID, passphrase: passphrase})
  }

  render() {
    if (this.state.page === Pages.GAME) {
      return <PokerTable
        gameID={this.state.gameID}
        initialPassphrase={this.state.passphrase}
        key={this.state.gameID}
      />
    } else if (this.state.page === Pages.CREATE_GAME) {
      return (
        <CreateGameComponent
          onGameCreateSuccess={this.handleCreateGameSuccess}
        />
      );
    } else if (this.state.page === Pages.MAIN) {
      return (
        <div style={{width: '400px', position: 'absolute'}} class="center">
          <button class="main-button" type="button" onClick={this.handleCreateClick}>Create Game</button>
          <text>Want to join a game? Just paste the link your friend sends you</text>
        </div>
      );
    } else {
      throw Error("Unknown page state " + this.state.page);
    }
  }
}

export default MainComponent;