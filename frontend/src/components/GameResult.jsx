import './GameResult.css'

function GameResult({ winner, player1Name, player2Name, username, isBot, onPlayAgain }) {
  const isDraw = winner === 'draw'
  const isVictory = winner === username
  const isBotVictory = winner === 'Bot'

  const getMessage = () => {
    if (isDraw) {
      return "It's a Draw! 🤝"
    }
    if (isVictory) {
      return 'You Won! 🎉'
    }
    return 'You Lost! 😢'
  }

  const getOpponentName = () => {
    return username === player1Name ? player2Name : player1Name
  }

  return (
    <div className="result-container">
      <div className={`result-card card ${isDraw ? 'draw' : isVictory ? 'victory' : 'defeat'}`}>
        <div className="result-header">
          <h1>{getMessage()}</h1>
        </div>

        <div className="result-details">
          {isDraw ? (
            <p>Both players played well!</p>
          ) : isVictory ? (
            <>
              <p className="victory-text">
                You defeated {isBotVictory ? 'the Bot' : getOpponentName()}!
              </p>
              <p className="victory-emoji">🏆</p>
            </>
          ) : (
            <p className="defeat-text">
              {isBotVictory ? 'The Bot won this round' : `${winner} won this round`}
            </p>
          )}
        </div>

        <div className="result-stats">
          <div className="stat">
            <p className="stat-label">Player 1</p>
            <p className="stat-value">{player1Name}</p>
          </div>
          <div className="vs-badge">VS</div>
          <div className="stat">
            <p className="stat-label">Player 2</p>
            <p className="stat-value">{player2Name}</p>
          </div>
        </div>

        <div className="result-actions">
          <button onClick={onPlayAgain} className="btn-play-again">
            Play Again
          </button>
        </div>
      </div>
    </div>
  )
}

export default GameResult
