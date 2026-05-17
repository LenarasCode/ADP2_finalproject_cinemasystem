// Tab switching
function showTab(tabName) {
  document.querySelectorAll('.tab-content').forEach(el => el.classList.remove('active'));
  document.getElementById(tabName).classList.add('active');
  document.querySelectorAll('.tab').forEach(btn => btn.classList.remove('active'));
  event.target.classList.add('active');
  if (tabName === 'notifications') fetchLogs();
  if (tabName === 'movies') fetchMovies();
}

// --- Movies ---
async function fetchMovies() {
  // Статический список 50 фильмов (позже можно заменить на gRPC вызов)
  const movies = [
    "The Shawshank Redemption|Drama|142","The Godfather|Crime|175","The Dark Knight|Action|152","Pulp Fiction|Crime|154","Schindler's List|History|195","Forrest Gump|Drama|142","Inception|Sci-Fi|148","Fight Club|Drama|139","The Matrix|Sci-Fi|136","Interstellar|Sci-Fi|169","Parasite|Thriller|132","Gladiator|Action|155","The Lion King|Animation|88","Avengers: Endgame|Action|181","Joker|Crime|122","Spirited Away|Animation|125","The Silence of the Lambs|Thriller|118","Se7en|Crime|127","The Prestige|Mystery|130","Whiplash|Drama|106","The Wolf of Wall Street|Comedy|180","Dune|Sci-Fi|155","The Truman Show|Comedy|103","Coco|Animation|105","Top Gun: Maverick|Action|130","La La Land|Musical|128","The Green Mile|Drama|189","Goodfellas|Crime|146","Saving Private Ryan|War|169","The Departed|Crime|151","American History X|Drama|119","Memento|Mystery|113","Eternal Sunshine of the Spotless Mind|Romance|108","The Grand Budapest Hotel|Comedy|99","Mad Max: Fury Road|Action|120","No Country for Old Men|Thriller|122","Oldboy|Mystery|120","Snatch|Comedy|102","The Intouchables|Comedy|112","A Beautiful Mind|Biography|135","Life Is Beautiful|Comedy|116","The Usual Suspects|Crime|106","Back to the Future|Adventure|116","Jurassic Park|Adventure|127","The Lord of the Rings: The Return of the King|Fantasy|201","Star Wars: Episode V|Sci-Fi|124","The Social Network|Biography|120","Good Will Hunting|Drama|126","The Shining|Horror|146","Alien|Horror|117"
  ];
  const tbody = document.getElementById('moviesBody');
  tbody.innerHTML = movies.map(m => {
    const [title, genre, dur] = m.split('|');
    return `<tr><td>${title}</td><td>${genre}</td><td>${dur} min</td></tr>`;
  }).join('');
}

// --- Notifications ---
async function fetchLogs() {
  const mockLogs = [
    { id: '1', recipient: 'customer@example.com', subject: 'Booking Confirmed', status: 'OK', sent_at: '2026-05-17 12:00' },
    { id: '2', recipient: 'vip@luxurycinema.com', subject: 'VIP Invitation', status: 'OK', sent_at: '2026-05-17 11:45' },
  ];
  renderLogs(mockLogs);
}
function renderLogs(logs) {
  document.getElementById('logsBody').innerHTML = logs.map(log => `
    <tr><td>${log.id}</td><td>${log.recipient}</td><td>${log.subject}</td><td>${log.status}</td><td>${log.sent_at}</td></tr>
  `).join('');
}
document.getElementById('sendBtn')?.addEventListener('click', () => {
  const email = document.getElementById('emailInput').value.trim();
  const booking = document.getElementById('bookingInput').value.trim();
  document.getElementById('statusMsg').textContent = email && booking ? `Email sent to ${email} (simulated).` : 'Fill all fields.';
});

// --- Booking ---
document.getElementById('bookBtn')?.addEventListener('click', () => {
  const userId = document.getElementById('userId').value.trim();
  const showtimeId = document.getElementById('showtimeId').value.trim();
  const seats = document.getElementById('seatIds').value.trim();
  if (!userId || !showtimeId || !seats) {
    document.getElementById('bookStatus').textContent = 'Fill all fields.';
    return;
  }
  document.getElementById('bookStatus').textContent = `Booking request sent for seats ${seats} (simulated).`;
  // В реальном проекте здесь POST к API Gateway
});

// Инициализация при загрузке
fetchMovies();
