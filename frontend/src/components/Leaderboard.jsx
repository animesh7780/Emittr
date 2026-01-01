import { useState, useEffect } from 'react'
import './Leaderboard.css'

function Leaderboard() {
  const [leaderboard, setLeaderboard] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchLeaderboard = async () => {
      try {
        const response = await fetch('/api/leaderboard')
        if (!response.ok) throw new Error('Failed to fetch leaderboard')
        const data = await response.json()
        setLeaderboard(data.leaderboard || [])
      } catch (err) {
        setError(err.message)
      } finally {
        setLoading(false)
      }
    }

    fetchLeaderboard()
  }, [])

  if (loading) {
    return (
      <div className="card text-center">
        <h2>Loading Leaderboard...</h2>
        <div className="spinner"></div>
      </div>
    )
  }

  return (
    <div className="leaderboard-container">
      <div className="card">
        <h2>🏆 Leaderboard</h2>

        {error && (
          <p className="error-message">{error}</p>
        )}

        {leaderboard.length === 0 ? (
          <p className="empty-message">No games played yet. Be the first!</p>
        ) : (
          <div className="leaderboard-table">
            <div className="leaderboard-header">
              <div className="rank-col">Rank</div>
              <div className="name-col">Player</div>
              <div className="stat-col">Wins</div>
              <div className="stat-col">Losses</div>
              <div className="stat-col">Draws</div>
              <div className="stat-col">Win Rate</div>
            </div>
            {leaderboard.map((player, idx) => (
              <div key={idx} className={`leaderboard-row ${idx < 3 ? `rank-${idx + 1}` : ''}`}>
                <div className="rank-col">
                  {idx === 0 && '🥇'}
                  {idx === 1 && '🥈'}
                  {idx === 2 && '🥉'}
                  {idx >= 3 && `#${idx + 1}`}
                </div>
                <div className="name-col">{player.username}</div>
                <div className="stat-col">{player.wins}</div>
                <div className="stat-col">{player.losses}</div>
                <div className="stat-col">{player.draws}</div>
                <div className="stat-col">{player.winRate}</div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default Leaderboard
