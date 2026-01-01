import { useState } from 'react'
import './Lobby.css'

function Lobby({ onRegister, onViewLeaderboard }) {
  const [username, setUsername] = useState('')

  const handleSubmit = (e) => {
    e.preventDefault()
    if (username.trim()) {
      onRegister(username.trim())
    }
  }

  return (
    <div className="lobby-container">
      <div className="card lobby-card">
        <h1>🎮 4 in a Row</h1>
        <p className="subtitle">Connect Four Discs to Win!</p>

        <form onSubmit={handleSubmit} className="lobby-form">
          <input
            type="text"
            placeholder="Enter your username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            maxLength={30}
            required
            autoFocus
          />
          <button type="submit" className="btn-primary">
            Play Now
          </button>
        </form>

        <button
          onClick={onViewLeaderboard}
          className="btn-secondary"
          style={{ marginTop: '12px' }}
        >
          View Leaderboard
        </button>

        <div className="rules">
          <h3>How to Play</h3>
          <ul>
            <li>Drop your disc into any column</li>
            <li>First to connect 4 discs wins!</li>
            <li>Play horizontally, vertically, or diagonally</li>
            <li>No opponent? Play against our bot!</li>
          </ul>
        </div>
      </div>
    </div>
  )
}

export default Lobby
