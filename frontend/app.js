const API_URL = window.location.origin;

// State
let jwtToken = localStorage.getItem('jwt');
let currentUserEmail = localStorage.getItem('email');
let isLoginMode = true;
let tickets = [];
let activeTicketId = null;

// DOM Elements
const authContainer = document.getElementById('auth-container');
const dashboardContainer = document.getElementById('dashboard-container');
const authForm = document.getElementById('auth-form');
const emailInput = document.getElementById('email');
const passwordInput = document.getElementById('password');
const authSubmitBtn = document.getElementById('auth-submit-btn');
const authToggleText = document.getElementById('auth-toggle-text');
const authToggleLink = document.getElementById('auth-toggle-link');
const authError = document.getElementById('auth-error');
const userEmailDisplay = document.getElementById('user-email-display');
const logoutBtn = document.getElementById('logout-btn');
const ticketsList = document.getElementById('tickets-list');
const newTicketBtn = document.getElementById('new-ticket-btn');

// Modals
const ticketModal = document.getElementById('ticket-modal');
const createTicketForm = document.getElementById('create-ticket-form');
const ticketTitleInput = document.getElementById('title');
const ticketDescInput = document.getElementById('description');
const cancelTicketBtn = document.getElementById('cancel-ticket-btn');
const ticketError = document.getElementById('ticket-error');

const viewModal = document.getElementById('view-modal');
const closeViewBtn = document.getElementById('close-view-btn');
const viewTitle = document.getElementById('view-title');
const viewStatusBadge = document.getElementById('view-status-badge');
const viewDate = document.getElementById('view-date');
const viewDesc = document.getElementById('view-desc');
const btnInProgress = document.getElementById('btn-in-progress');
const btnClosed = document.getElementById('btn-closed');
const updateError = document.getElementById('update-error');

// Initialization
function init() {
    if (jwtToken) {
        showDashboard();
    } else {
        showAuth();
    }
}

// UI Toggles
function showAuth() {
    authContainer.classList.add('active');
    dashboardContainer.classList.remove('active');
}

function showDashboard() {
    authContainer.classList.remove('active');
    dashboardContainer.classList.add('active');
    userEmailDisplay.textContent = currentUserEmail;
    fetchTickets();
}

authToggleLink.addEventListener('click', (e) => {
    e.preventDefault();
    isLoginMode = !isLoginMode;
    authSubmitBtn.textContent = isLoginMode ? 'Login' : 'Register';
    authToggleText.textContent = isLoginMode ? "Don't have an account?" : "Already have an account?";
    authToggleLink.textContent = isLoginMode ? 'Register' : 'Login';
    authError.classList.add('hidden');
});

function showError(el, msg) {
    el.textContent = msg;
    el.classList.remove('hidden');
}

// API Calls
async function apiRequest(endpoint, method = 'GET', body = null) {
    const headers = {
        'Content-Type': 'application/json'
    };
    if (jwtToken) {
        headers['Authorization'] = `Bearer ${jwtToken}`;
    }

    const config = { method, headers };
    if (body) config.body = JSON.stringify(body);

    const res = await fetch(`${API_URL}${endpoint}`, config);
    const data = await res.json().catch(() => ({}));

    if (!res.ok) {
        throw new Error(data.error || 'Something went wrong');
    }
    return data;
}

// Auth Handlers
authForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    authError.classList.add('hidden');
    const email = emailInput.value;
    const password = passwordInput.value;
    
    try {
        if (!isLoginMode) {
            await apiRequest('/auth/register', 'POST', { email, password });
            // Auto login after register
        }
        const data = await apiRequest('/auth/login', 'POST', { email, password });
        jwtToken = data.token;
        currentUserEmail = email;
        localStorage.setItem('jwt', jwtToken);
        localStorage.setItem('email', email);
        showDashboard();
    } catch (err) {
        showError(authError, err.message);
    }
});

logoutBtn.addEventListener('click', () => {
    jwtToken = null;
    currentUserEmail = null;
    localStorage.removeItem('jwt');
    localStorage.removeItem('email');
    showAuth();
});

// Ticket Handlers
async function fetchTickets() {
    try {
        const data = await apiRequest('/tickets');
        tickets = data.tickets || [];
        renderTickets();
    } catch (err) {
        if (err.message.includes('token') || err.message.includes('unauthorized')) {
            logoutBtn.click();
        }
    }
}

function renderTickets() {
    if (tickets.length === 0) {
        ticketsList.innerHTML = `<p style="grid-column: 1/-1; text-align: center; color: var(--text-secondary); padding: 2rem;">No tickets found. Create your first one!</p>`;
        return;
    }

    ticketsList.innerHTML = tickets.map(t => `
        <div class="ticket-card" onclick="openTicket(${t.id})">
            <div class="ticket-header">
                <h4>${t.title}</h4>
                <span class="status-badge status-${t.status}">${t.status.replace('_', ' ')}</span>
            </div>
            <p class="ticket-desc">${t.description}</p>
            <p class="date-text">${new Date(t.created_at).toLocaleDateString()}</p>
        </div>
    `).join('');
}

// Create Ticket Flow
newTicketBtn.addEventListener('click', () => {
    ticketModal.classList.remove('hidden');
    ticketError.classList.add('hidden');
    createTicketForm.reset();
});

cancelTicketBtn.addEventListener('click', () => {
    ticketModal.classList.add('hidden');
});

createTicketForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    ticketError.classList.add('hidden');
    
    try {
        await apiRequest('/tickets', 'POST', {
            title: ticketTitleInput.value,
            description: ticketDescInput.value
        });
        ticketModal.classList.add('hidden');
        fetchTickets();
    } catch (err) {
        showError(ticketError, err.message);
    }
});

// View/Update Ticket Flow
async function openTicket(id) {
    activeTicketId = id;
    try {
        const t = await apiRequest(`/tickets/${id}`);
        viewTitle.textContent = t.title;
        viewDesc.textContent = t.description;
        viewDate.textContent = `Created at: ${new Date(t.created_at).toLocaleString()}`;
        
        viewStatusBadge.textContent = t.status.replace('_', ' ');
        viewStatusBadge.className = `status-badge status-${t.status}`;

        // Configure buttons based on strict flow: open -> in_progress -> closed
        btnInProgress.classList.add('hidden');
        btnClosed.classList.add('hidden');
        updateError.classList.add('hidden');

        if (t.status === 'open') {
            btnInProgress.classList.remove('hidden');
        } else if (t.status === 'in_progress') {
            btnClosed.classList.remove('hidden');
        }

        viewModal.classList.remove('hidden');
    } catch (err) {
        alert(err.message);
    }
}

closeViewBtn.addEventListener('click', () => {
    viewModal.classList.add('hidden');
    activeTicketId = null;
});

async function updateStatus(newStatus) {
    if (!activeTicketId) return;
    try {
        await apiRequest(`/tickets/${activeTicketId}/status`, 'PATCH', { status: newStatus });
        openTicket(activeTicketId); // Reload modal data
        fetchTickets(); // Update background list
    } catch (err) {
        showError(updateError, err.message);
    }
}

btnInProgress.addEventListener('click', () => updateStatus('in_progress'));
btnClosed.addEventListener('click', () => updateStatus('closed'));

// Start app
init();
