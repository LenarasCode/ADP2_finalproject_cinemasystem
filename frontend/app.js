document.addEventListener('DOMContentLoaded', () => {
  fetchLogs();
  document.getElementById('sendBtn').addEventListener('click', sendEmail);
});

async function fetchLogs() {
  // TODO: заменить на реальный запрос к API Gateway
  const mockLogs = [
    { id: '1', recipient: 'customer@example.com', subject: 'Booking Confirmed', status: 'OK', sent_at: '2026-05-17 12:00' },
    { id: '2', recipient: 'vip@luxurycinema.com', subject: 'VIP Invitation', status: 'OK', sent_at: '2026-05-17 11:45' },
  ];
  renderLogs(mockLogs);
}

function renderLogs(logs) {
  const tbody = document.getElementById('logsBody');
  tbody.innerHTML = logs.map(log => `
    <tr>
      <td>${log.id}</td>
      <td>${log.recipient}</td>
      <td>${log.subject}</td>
      <td><span class="badge badge-ok">${log.status}</span></td>
      <td>${log.sent_at}</td>
    </tr>
  `).join('');
}

async function sendEmail() {
  const email = document.getElementById('emailInput').value.trim();
  const booking = document.getElementById('bookingInput').value.trim();
  const status = document.getElementById('statusMsg');
  if (!email || !booking) {
    status.textContent = 'Please fill both fields.';
    return;
  }
  // TODO: POST к API Gateway
  status.textContent = `Email sent to ${email} for booking ${booking} (simulated).`;
  document.getElementById('emailInput').value = '';
  document.getElementById('bookingInput').value = '';
  setTimeout(fetchLogs, 1000);
}
