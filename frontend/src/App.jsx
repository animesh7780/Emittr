import { useState, useEffect, useCallback } from 'react'
import GameBoard from './components/GameBoard'
import Lobby from './components/Lobby'
import Leaderboard from './components/Leaderboard'
import GameResult from './components/GameResult'
import './App.css'

function App() {
  const [gameState, setGameState] = useState('lobby') // 'lobby', 'waiting', 'playing', 'result'
  const [username, setUsername] = useState('')
  const [gameId, setGameId] = useState('')
  const [board, setBoard] = useState(Array(6).fill(null).map(() => Array(7).fill(0)))
  const [currentPlayer, setCurrentPlayer] = useState(1)
  const [isBot, setIsBot] = useState(false)
  const [player1Name, setPlayer1Name] = useState('')
  const [player2Name, setPlayer2Name] = useState('')
  const [winner, setWinner] = useState('')
  const [winRow, setWinRow] = useState(-1)
  const [winCol, setWinCol] = useState(-1)
  const [ws, setWs] = useState(null)
  const [error, setError] = useState('')
  const [showLeaderboard, setShowLeaderboard] = useState(false)

  // WebSocket connection
  useEffect(() => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    // Connect directly to backend for WebSocket (bypasses nginx)
    const wsUrl = `${protocol}//emittr-backend-0l70.onrender.com/ws`
    const websocket = new WebSocket(wsUrl)

    websocket.onopen = () => {
      console.log('WebSocket connected')
      setError('')
    }

    websocket.onmessage = (event) => {
      const data = JSON.parse(event.data)
      console.log('Message received:', data)

      switch (data.type) {
        case 'game_start':
          setGameId(data.payload.gameId)
          setPlayer1Name(data.payload.player1)
          setPlayer2Name(data.payload.player2)
          setIsBot(data.payload.isBot)
          setCurrentPlayer(data.payload.yourTurn ? 1 : 2)
          setBoard(Array(6).fill(null).map(() => Array(7).fill(0)))
          setWinner('')
          setGameState('playing')
          break
        case 'game_move':
          setBoard(data.payload.board)
          setCurrentPlayer(data.payload.currentPlayer)
          break
        case 'game_result':
          setWinner(data.payload.winner)
          setWinRow(data.payload.winRow || -1)
          setWinCol(data.payload.winCol || -1)
          setGameState('result')
          break
        case 'error':
          setError(data.payload.message)
          break
        default:
          console.log('Unknown message type:', data.type)
      }
    }

    websocket.onerror = (error) => {
      console.error('WebSocket error:', error)
      setError('Connection error. Please refresh and try again.')
    }

    websocket.onclose = () => {
      console.log('WebSocket disconnected')
    }

    setWs(websocket)

    return () => {
      websocket.close()
    }
  }, [])

  const handleRegister = useCallback((name) => {
    setUsername(name)
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'register',
        payload: { username: name }
      }))
      setGameState('waiting')
    }
  }, [ws])

  const handleGameStart = useCallback((payload) => {
    setGameId(payload.gameId)
    setPlayer1Name(payload.player1)
    setPlayer2Name(payload.player2)
    setIsBot(payload.isBot)
    setCurrentPlayer(payload.yourTurn ? 1 : 2)
    setBoard(Array(6).fill(null).map(() => Array(7).fill(0)))
    setWinner('')
    setGameState('playing')
  }, [])

  const handleGameMove = useCallback((payload) => {
    const newBoard = payload.board
    setBoard(newBoard)
    setCurrentPlayer(payload.player === 1 ? 2 : 1)
  }, [])

  const handleGameResult = useCallback((payload) => {
    setWinner(payload.winner)
    setWinRow(payload.winRow || -1)
    setWinCol(payload.winCol || -1)
    setGameState('result')
  }, [])

  const handleColumnClick = useCallback((col) => {
    if (gameState === 'playing' && ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'game_move',
        payload: { column: col }
      }))
    }
  }, [gameState, ws])

  const handlePlayAgain = useCallback(() => {
    setGameState('lobby')
    setUsername('')
    setGameId('')
    setBoard(Array(6).fill(null).map(() => Array(7).fill(0)))
    setCurrentPlayer(1)
    setIsBot(false)
    setPlayer1Name('')
    setPlayer2Name('')
    setWinner('')
    setError('')
  }, [])

  return (
    <div className="container">
      {error && (
        <div className="error-banner">
          <p>{error}</p>
          <button onClick={() => setError('')}>Dismiss</button>
        </div>
      )}

      {showLeaderboard ? (
        <>
          <Leaderboard />
          <div className="button-group">
            <button
              className="btn-secondary"
              onClick={() => setShowLeaderboard(false)}
            >
              Back to Game
            </button>
          </div>
        </>
      ) : (
        <>
          {gameState === 'lobby' && (
            <Lobby onRegister={handleRegister} onViewLeaderboard={() => setShowLeaderboard(true)} />
          )}

          {gameState === 'waiting' && (
            <div className="card text-center">
              <h2>Waiting for opponent...</h2>
              <p>Finding a player or starting with bot in 10 seconds</p>
              <div className="spinner"></div>
            </div>
          )}

          {gameState === 'playing' && (
            <GameBoard
              board={board}
              currentPlayer={currentPlayer}
              player1Name={player1Name}
              player2Name={player2Name}
              onColumnClick={handleColumnClick}
              username={username}
              isBot={isBot}
              winRow={winRow}
              winCol={winCol}
            />
          )}

          {gameState === 'result' && (
            <GameResult
              winner={winner}
              player1Name={player1Name}
              player2Name={player2Name}
              username={username}
              isBot={isBot}
              onPlayAgain={handlePlayAgain}
            />
          )}
        </>
      )}
    </div>
  )
}

export default App
