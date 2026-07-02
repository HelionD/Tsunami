// ── Tab switch ──
function switchTab(tab) {
    document.querySelectorAll('.auth-tab').forEach((t, i) => {
        t.classList.toggle('active', (tab === 'signin' && i === 0) || (tab === 'signup' && i === 1));
    });
    document.getElementById('signinForm').classList.toggle('active', tab === 'signin');
    document.getElementById('signupForm').classList.toggle('active', tab === 'signup');
}

// ── Sign in / sign up → show dashboard ──
function handleSignIn() {
    // Get name if signing up
    const first = document.getElementById('su-first')?.value.trim() ||
                  document.getElementById('si-email')?.value.split('@')[0] || 'Engineer';
    document.getElementById('dashUserName').textContent = first.charAt(0).toUpperCase() + first.slice(1);

    document.getElementById('authScreen').style.display = 'none';
    document.getElementById('dashboardScreen').style.display = 'block';
}

// ── Log out → back to auth ──
function handleLogout() {
    document.getElementById('dashboardScreen').style.display = 'none';
    document.getElementById('authScreen').style.display = 'flex';
    // Clear inputs
    document.querySelectorAll('.auth-form input').forEach(i => i.value = '');
    switchTab('signin');
}
