import './GameBoard.css'

function GameBoard({ board, currentPlayer, player1Name, player2Name, onColumnClick, username, isBot, winRow, winCol }) {
  const isYourTurn = (currentPlayer === 1 && username === player1Name) || (currentPlayer === 2 && username === player2Name)

  const renderCell = (row, col) => {
    const value = board[row][col]
    const isWinCell = winRow === row && winCol === col
    
    let cellClass = 'cell'
    if (value === 1) cellClass += ' player1'
    if (value === 2) cellClass += ' player2'
    if (isWinCell) cellClass += ' win-cell'

    return (
      <div key={`${row}-${col}`} className={cellClass}>
        {value !== 0 && <div className="disc"></div>}
      </div>
    )
  }

  return (
    <div className="game-container">
      <div className="game-info card">
        <div className="player-info">
          <div className={`player ${currentPlayer === 1 ? 'active' : ''}`}>
            <div className={`player-disc player1-disc`}></div>
            <div>
              <p className="player-name">{player1Name}</p>
              <p className="player-label">Player 1</p>
            </div>
          </div>
          <div className="vs">VS</div>
          <div className={`player ${currentPlayer === 2 ? 'active' : ''}`}>
            <div>
              <p className="player-name">{player2Name}</p>
              <p className="player-label">{isBot ? 'Bot' : 'Player 2'}</p>
            </div>
            <div className={`player-disc player2-disc`}></div>
          </div>
        </div>

        <div className="turn-indicator">
          {isYourTurn ? (
            <p className="your-turn">🎯 Your Turn!</p>
          ) : (
            <p className="opponent-turn">⏳ Waiting for opponent...</p>
          )}
        </div>
      </div>

      <div className="board-container card">
        <div className="board">
          {board.map((row, rowIdx) =>
            row.map((_, colIdx) => renderCell(rowIdx, colIdx))
          )}
        </div>

        <div className="column-buttons">
          {[0, 1, 2, 3, 4, 5, 6].map((col) => (
            <button
              key={col}
              className="column-btn"
              onClick={() => onColumnClick(col)}
              disabled={!isYourTurn}
              title={isYourTurn ? `Drop in column ${col}` : 'Not your turn'}
            >
              ↓
            </button>
          ))}
        </div>
      </div>
    </div>
  )
}

export default GameBoard
